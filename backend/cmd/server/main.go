package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"mastercard-backend/internal/config"
	"mastercard-backend/internal/database"
	"mastercard-backend/internal/handlers"
	"mastercard-backend/internal/middleware"
	"mastercard-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize query service (Gemini client)
	queryService, err := services.NewQueryService()
	if err != nil {
		log.Fatalf("Failed to initialize query service: %v", err)
	}
	defer queryService.Close()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      config.AppConfig.AppName,
		ErrorHandler: errorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(middleware.LoggerMiddleware())
	app.Use(middleware.RequestIDMiddleware())
	app.Use(middleware.CORSMiddleware())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": config.AppConfig.AppName,
		})
	})

	// Initialize handlers
	authHandler := handlers.NewAuthHandler()
	queryHandler := handlers.NewQueryHandler(queryService)
	conversationHandler := handlers.NewConversationHandler()
	adminHandler := handlers.NewAdminHandler()

	// Public routes
	api := app.Group("/api/v1")
	{
		// Authentication routes
		auth := api.Group("/auth")
		{
			auth.Post("/register", authHandler.Register)
			auth.Post("/login", authHandler.Login)
			auth.Post("/refresh", authHandler.RefreshToken)
		}
	}

	// Protected routes
	protected := api.Group("", middleware.AuthMiddleware())
	{
		// User profile
		protected.Get("/auth/profile", authHandler.GetProfile)

		// Query routes
		queries := protected.Group("/query")
		{
			queries.Post("", queryHandler.ExecuteQuery)
		}

		// Conversation routes
		conversations := protected.Group("/conversations")
		{
			conversations.Post("", conversationHandler.CreateConversation)
			conversations.Get("", conversationHandler.GetConversations)
			conversations.Get("/search", conversationHandler.SearchConversations)
			conversations.Get("/:id", conversationHandler.GetConversation)
			conversations.Put("/:id", conversationHandler.UpdateConversation)
			conversations.Delete("/:id", conversationHandler.DeleteConversation)
			conversations.Post("/:id/branch", conversationHandler.CreateBranch)
		}

		// Admin routes (Manager and Admin access)
		admin := protected.Group("/admin", middleware.RequireRole("manager", "admin"))
		{
			// User management (Manager and Admin can create, read, update)
			admin.Get("/users", adminHandler.GetUsers)
			admin.Post("/users", adminHandler.CreateUser)
			admin.Put("/users/:id", adminHandler.UpdateUser)
			
			// Audit logs (Manager and Admin can view)
			admin.Get("/audit-logs", adminHandler.GetAuditLogs)
			
			// Metrics (Manager and Admin can view)
			admin.Get("/metrics", adminHandler.GetMetrics)
			
			// User deletion (Admin only)
			admin.Delete("/users/:id", middleware.RequireRole("admin"), adminHandler.DeleteUser)
		}
	}

	// Start server
	addr := config.AppConfig.AppHost + ":" + config.AppConfig.AppPort
	log.Printf("Server starting on %s", addr)

	// Graceful shutdown
	go func() {
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited")
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"error": message,
	})
}

