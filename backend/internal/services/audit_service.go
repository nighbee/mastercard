package services

import (
	"mastercard-backend/internal/database"
	"mastercard-backend/internal/models"
	"time"
)

type AuditService struct{}

func NewAuditService() *AuditService {
	return &AuditService{}
}

// LogAction logs an action to the audit log
func (s *AuditService) LogAction(userID *uint, action, resource string, queryText, sqlExecuted *string, resultCount *int, ipAddress, userAgent *string, status string, errorMessage *string, executionTimeMs *int) error {
	auditLog := models.AuditLog{
		UserID:          userID,
		Action:          action,
		Resource:        &resource,
		QueryText:       queryText,
		SQLExecuted:     sqlExecuted,
		ResultCount:     resultCount,
		IPAddress:       ipAddress,
		UserAgent:       userAgent,
		Status:          &status,
		ErrorMessage:    errorMessage,
		ExecutionTimeMs: executionTimeMs,
		Timestamp:       time.Now(),
	}

	return database.DB.Create(&auditLog).Error
}

// GetAuditLogs retrieves audit logs with filtering
func (s *AuditService) GetAuditLogs(userID *uint, limit, offset int, action, resource, status *string) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	query := database.DB.Model(&models.AuditLog{})

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	if action != nil {
		query = query.Where("action = ?", *action)
	}

	if resource != nil {
		query = query.Where("resource = ?", *resource)
	}

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get logs
	if err := query.Order("timestamp DESC").
		Limit(limit).
		Offset(offset).
		Preload("User").
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

