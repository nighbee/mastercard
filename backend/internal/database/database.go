package database

import (
	"fmt"
	"log"

	"mastercard-backend/internal/config"
	"mastercard-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect initializes the database connection
func Connect() error {
	// --- START DEBUGGING ---
	// Print the configuration to see exactly what the application is using.
	log.Println("--- DATABASE CONNECTION CONFIG ---")
	log.Printf("HOST: [%s]", config.AppConfig.DBHost)
	log.Printf("PORT: [%s]", config.AppConfig.DBPort)
	log.Printf("USER: [%s]", config.AppConfig.DBUser)
	log.Printf("PASSWORD: [%s]", config.AppConfig.DBPassword) // This will reveal the actual password being used.
	log.Printf("DBNAME: [%s]", config.AppConfig.DBName)
	log.Printf("SSLMODE: [%s]", config.AppConfig.DBSSLMode)
	log.Println("---------------------------------")
	// --- END DEBUGGING ---

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.AppConfig.DBHost,
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBName,
		config.AppConfig.DBPort,
		config.AppConfig.DBSSLMode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		// The original error from gorm is very informative, so we wrap it.
		return fmt.Errorf("failed to initialize database, got error %w", err)
	}

	log.Println("Database connected successfully")

	// Auto-migration is disabled. Schema changes should be handled by SQL migration files
	// in the /migrations directory. This prevents conflicts with manually defined policies and types.
	// if err := autoMigrate(); err != nil {
	// 	return fmt.Errorf("failed to auto-migrate: %w", err)
	// }

	return nil
}

// autoMigrate runs GORM auto-migration
func autoMigrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.Transaction{},
		&models.Conversation{},
		&models.Message{},
		&models.AuditLog{},
	)
}

// SetCurrentUserID sets the current user ID for RLS policies
// This should be called before executing queries that need RLS
func SetCurrentUserID(userID uint) error {
	return DB.Exec("SET app.current_user_id = ?", userID).Error
}

// ClearCurrentUserID clears the current user ID setting
func ClearCurrentUserID() error {
	return DB.Exec("RESET app.current_user_id").Error
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
