package services

import (
	"errors"
	"time"

	"mastercard-backend/internal/database"
	"mastercard-backend/internal/models"
	"mastercard-backend/internal/utils"

	"gorm.io/gorm"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

// Register creates a new user account
func (s *AuthService) Register(email, password, fullName string, roleID *uint) (*models.User, error) {
	// Check if user already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := models.User{
		Email:        email,
		PasswordHash: hashedPassword,
		FullName:     fullName,
		RoleID:       roleID,
		IsActive:     true,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return nil, errors.New("failed to create user")
	}

	// Load role
	if user.RoleID != nil {
		database.DB.Preload("Role").First(&user, user.ID)
	}

	return &user, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(email, password string) (*models.User, string, string, error) {
	var user models.User
	if err := database.DB.Preload("Role").Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, "", "", errors.New("invalid email or password")
		}
		return nil, "", "", errors.New("database error")
	}

	if !user.IsActive {
		return nil, "", "", errors.New("user account is inactive")
	}

	// Check password
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, "", "", errors.New("invalid email or password")
	}

	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, user.RoleID)
	if err != nil {
		return nil, "", "", errors.New("failed to generate access token")
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, "", "", errors.New("failed to generate refresh token")
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	database.DB.Save(&user)

	return &user, accessToken, refreshToken, nil
}

// RefreshToken generates new access token from refresh token
func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	claims, err := utils.ValidateToken(refreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	var user models.User
	if err := database.DB.Preload("Role").First(&user, claims.UserID).Error; err != nil {
		return "", errors.New("user not found")
	}

	if !user.IsActive {
		return "", errors.New("user account is inactive")
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, user.RoleID)
	if err != nil {
		return "", errors.New("failed to generate access token")
	}

	return accessToken, nil
}

