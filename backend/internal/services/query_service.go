package services

import (
	"context"
	"encoding/json"

	// "errors"
	"fmt"
	"strings"
	"time"

	"mastercard-backend/internal/config"
	"mastercard-backend/internal/database"
	"mastercard-backend/internal/models"
	"mastercard-backend/pkg/gemini"
)

type QueryService struct {
	geminiClient *gemini.Client
}

func NewQueryService() (*QueryService, error) {
	client, err := gemini.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Gemini client: %w", err)
	}

	return &QueryService{
		geminiClient: client,
	}, nil
}

// ExecuteQuery processes a natural language query and returns results
func (s *QueryService) ExecuteQuery(userID uint, query string, conversationID *uint) (*models.Message, error) {
	startTime := time.Now()

	// Get conversation history if conversationID is provided
	var history []string
	if conversationID != nil && *conversationID > 0 {
		var messages []models.Message
		database.DB.Where("conversation_id = ?", *conversationID).
			Order("created_at DESC").
			Limit(10).
			Find(&messages)

		// Reverse to get chronological order
		for i := len(messages) - 1; i >= 0; i-- {
			history = append(history, messages[i].UserMessage)
		}
	}

	// Generate SQL using Gemini
	schemaContext := gemini.GetSchemaContext()
	sqlQuery, err := s.geminiClient.GenerateSQL(query, schemaContext, history)
	if err != nil {
		return s.createErrorMessage(userID, conversationID, query, fmt.Sprintf("Failed to generate SQL: %v", err), startTime)
	}

	// Validate SQL (basic check - no DROP, DELETE, UPDATE, INSERT, TRUNCATE)
	if !s.isValidReadOnlyQuery(sqlQuery) {
		return s.createErrorMessage(userID, conversationID, query, "Only SELECT queries are allowed", startTime)
	}

	// Execute SQL query
	result, resultFormat, err := s.executeSQL(sqlQuery)
	executionTime := int(time.Since(startTime).Milliseconds())

	if err != nil {
		return s.createErrorMessage(userID, conversationID, query, fmt.Sprintf("Query execution failed: %v", err), startTime)
	}

	// Create message record
	message := models.Message{
		UserMessage:     query,
		SQLQuery:        &sqlQuery,
		ResultFormat:    &resultFormat,
		ExecutionTimeMs: &executionTime,
	}

	// Set result_data only if not empty (JSONB requires valid JSON or NULL)
	if result != "" {
		message.ResultData = &result
	}

	// Set conversation ID if provided
	if conversationID != nil && *conversationID > 0 {
		message.ConversationID = *conversationID
	}

	if err := database.DB.Create(&message).Error; err != nil {
		return nil, fmt.Errorf("failed to save message: %w", err)
	}

	return &message, nil
}

// executeSQL executes a SQL query and returns results
func (s *QueryService) executeSQL(sqlQuery string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.AppConfig.QueryTimeoutSeconds)*time.Second)
	defer cancel()

	// Get raw database connection
	sqlDB, err := database.DB.DB()
	if err != nil {
		return "", "", fmt.Errorf("failed to get database connection: %w", err)
	}

	// Execute query with timeout
	rows, err := sqlDB.QueryContext(ctx, sqlQuery)
	if err != nil {
		return "", "", fmt.Errorf("query execution error: %w", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return "", "", fmt.Errorf("failed to get columns: %w", err)
	}

	// Scan results
	var results []map[string]interface{}
	rowCount := 0
	maxRows := config.AppConfig.MaxResultRows

	for rows.Next() {
		if rowCount >= maxRows {
			break
		}

		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return "", "", fmt.Errorf("failed to scan row: %w", err)
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if val != nil {
				// Handle different types
				switch v := val.(type) {
				case []byte:
					row[col] = string(v)
				case time.Time:
					row[col] = v.Format(time.RFC3339)
				default:
					row[col] = v
				}
			} else {
				row[col] = nil
			}
		}
		results = append(results, row)
		rowCount++
	}

	if err := rows.Err(); err != nil {
		return "", "", fmt.Errorf("row iteration error: %w", err)
	}

	// Determine result format
	resultFormat := "table"
	if len(results) == 0 {
		resultFormat = "text"
	} else if len(results) == 1 && len(columns) == 1 {
		resultFormat = "text"
	}

	// Convert to JSON
	resultJSON, err := json.Marshal(results)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal results: %w", err)
	}

	return string(resultJSON), resultFormat, nil
}

// isValidReadOnlyQuery checks if the query is a safe read-only query
func (s *QueryService) isValidReadOnlyQuery(sql string) bool {
	upperSQL := strings.ToUpper(" " + sql + " ")
	dangerousKeywords := []string{" DROP ", " DELETE ", " UPDATE ", " INSERT ", " TRUNCATE ", " ALTER ", " CREATE ", " GRANT ", " REVOKE "}

	for _, keyword := range dangerousKeywords {
		if strings.Contains(upperSQL, keyword) {
			return false
		}
	}

	return strings.Contains(upperSQL, " SELECT ")
}

// createErrorMessage creates an error message record
func (s *QueryService) createErrorMessage(userID uint, conversationID *uint, query, errorMsg string, startTime time.Time) (*models.Message, error) {
	executionTime := int(time.Since(startTime).Milliseconds())
	format := "error"

	message := models.Message{
		UserMessage:     query,
		ErrorMessage:    &errorMsg,
		ResultFormat:    &format,
		ExecutionTimeMs: &executionTime,
	}

	// Set conversation ID if provided
	if conversationID != nil && *conversationID > 0 {
		message.ConversationID = *conversationID
	}

	if err := database.DB.Create(&message).Error; err != nil {
		return nil, fmt.Errorf("failed to save error message: %w", err)
	}

	return &message, nil
}

// Close closes the Gemini client
func (s *QueryService) Close() error {
	if s.geminiClient != nil {
		return s.geminiClient.Close()
	}
	return nil
}
