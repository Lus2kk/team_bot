package model

import (
	"time"
)

// LogLevel represents the level of the log entry
type LogLevel string

const (
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelError LogLevel = "ERROR"
	LogLevelDebug LogLevel = "DEBUG"
)

// OperationType represents the type of operation being logged
type OperationType string

const (
	OperationUserRegistration OperationType = "USER_REGISTRATION"
	OperationUserLogin        OperationType = "USER_LOGIN"
	OperationAdminAction      OperationType = "ADMIN_ACTION"
	OperationTokenGeneration  OperationType = "TOKEN_GENERATION"
	OperationTokenUsage       OperationType = "TOKEN_USAGE"
	OperationUserUpdate       OperationType = "USER_UPDATE"
	OperationBotCommand       OperationType = "BOT_COMMAND"
	OperationError            OperationType = "ERROR"
)

// OperationLog represents a log entry for an operation
type OperationLog struct {
	ID            int64         `json:"id"`
	UserID        *int64        `json:"user_id,omitempty"`
	ChatID        *int64        `json:"chat_id,omitempty"`
	Username      *string       `json:"username,omitempty"`
	OperationType OperationType `json:"operation_type"`
	Level         LogLevel      `json:"level"`
	Message       string        `json:"message"`
	Details       *string       `json:"details,omitempty"`
	IPAddress     *string       `json:"ip_address,omitempty"`
	UserAgent     *string       `json:"user_agent,omitempty"`
	Success       bool          `json:"success"`
	Duration      *int64        `json:"duration_ms,omitempty"` // Duration in milliseconds
	ErrorCode     *string       `json:"error_code,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
}

// LogContext contains contextual information for logging
type LogContext struct {
	UserID    *int64
	ChatID    *int64
	Username  *string
	IPAddress *string
	UserAgent *string
}


type LogEntry struct {
	Level         LogLevel
	OperationType OperationType
	Message       string
	Details       *string
	Context       *LogContext
	Success       bool
	Duration      *int64
	ErrorCode     *string
	Messages *int
}
