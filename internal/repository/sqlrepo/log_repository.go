package sqlrepo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"team_bot/internal/model"
)

type LogRepository struct {
	db *sql.DB
}

func NewLogRepository(db *sql.DB) *LogRepository {
	return &LogRepository{db: db}
}


func (r *LogRepository) LogOperation(ctx context.Context, entry model.LogEntry) error {
	query := `
		INSERT INTO operation_logs 
		(user_id, chat_id, username, operation_type, level, message, details, 
		 ip_address, user_agent, success, duration_ms, error_code, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	var userID, chatID *int64
	var username, ipAddress, userAgent *string

	if entry.Context != nil {
		userID = entry.Context.UserID
		chatID = entry.Context.ChatID
		username = entry.Context.Username
		ipAddress = entry.Context.IPAddress
		userAgent = entry.Context.UserAgent
	}

	_, err := r.db.ExecContext(ctx, query,
		userID,
		chatID,
		username,
		entry.OperationType,
		entry.Level,
		entry.Message,
		entry.Details,
		ipAddress,
		userAgent,
		entry.Success,
		entry.Duration,
		entry.ErrorCode,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to log operation: %v", err)
	}

	return nil
}


func (r *LogRepository) GetOperationLogs(ctx context.Context, filters LogFilters, limit, offset int) ([]*model.OperationLog, error) {
	query := `
		SELECT id, user_id, chat_id, username, operation_type, level, message, details,
		       ip_address, user_agent, success, duration_ms, error_code, created_at
		FROM operation_logs
		WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1


	if filters.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, *filters.UserID)
		argIndex++
	}

	if filters.ChatID != nil {
		query += fmt.Sprintf(" AND chat_id = $%d", argIndex)
		args = append(args, *filters.ChatID)
		argIndex++
	}

	if filters.OperationType != nil {
		query += fmt.Sprintf(" AND operation_type = $%d", argIndex)
		args = append(args, *filters.OperationType)
		argIndex++
	}

	if filters.Level != nil {
		query += fmt.Sprintf(" AND level = $%d", argIndex)
		args = append(args, *filters.Level)
		argIndex++
	}

	if filters.Success != nil {
		query += fmt.Sprintf(" AND success = $%d", argIndex)
		args = append(args, *filters.Success)
		argIndex++
	}

	if filters.StartTime != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", argIndex)
		args = append(args, *filters.StartTime)
		argIndex++
	}

	if filters.EndTime != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", argIndex)
		args = append(args, *filters.EndTime)
		argIndex++
	}


	query += " ORDER BY created_at DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query operation logs: %v", err)
	}
	defer rows.Close()

	var logs []*model.OperationLog
	for rows.Next() {
		log := &model.OperationLog{}
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.ChatID,
			&log.Username,
			&log.OperationType,
			&log.Level,
			&log.Message,
			&log.Details,
			&log.IPAddress,
			&log.UserAgent,
			&log.Success,
			&log.Duration,
			&log.ErrorCode,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan operation log: %v", err)
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over operation logs: %v", err)
	}

	return logs, nil
}


func (r *LogRepository) GetOperationLogStats(ctx context.Context, filters LogFilters) (*LogStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_count,
			COUNT(CASE WHEN success = true THEN 1 END) as success_count,
			COUNT(CASE WHEN success = false THEN 1 END) as error_count,
			AVG(CASE WHEN duration_ms IS NOT NULL THEN duration_ms END) as avg_duration_ms
		FROM operation_logs
		WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1


	if filters.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, *filters.UserID)
		argIndex++
	}

	if filters.ChatID != nil {
		query += fmt.Sprintf(" AND chat_id = $%d", argIndex)
		args = append(args, *filters.ChatID)
		argIndex++
	}

	if filters.OperationType != nil {
		query += fmt.Sprintf(" AND operation_type = $%d", argIndex)
		args = append(args, *filters.OperationType)
		argIndex++
	}

	if filters.Level != nil {
		query += fmt.Sprintf(" AND level = $%d", argIndex)
		args = append(args, *filters.Level)
		argIndex++
	}

	if filters.StartTime != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", argIndex)
		args = append(args, *filters.StartTime)
		argIndex++
	}

	if filters.EndTime != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", argIndex)
		args = append(args, *filters.EndTime)
		argIndex++
	}

	var stats LogStats
	var avgDuration sql.NullFloat64

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&stats.TotalCount,
		&stats.SuccessCount,
		&stats.ErrorCount,
		&avgDuration,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get operation log stats: %v", err)
	}

	if avgDuration.Valid {
		stats.AvgDurationMs = &avgDuration.Float64
	}

	return &stats, nil
}


func (r *LogRepository) DeleteOldLogs(ctx context.Context, olderThan time.Duration) (int64, error) {
	query := `DELETE FROM operation_logs WHERE created_at < $1`

	cutoffTime := time.Now().Add(-olderThan)
	result, err := r.db.ExecContext(ctx, query, cutoffTime)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old logs: %v", err)
	}

	deleted, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %v", err)
	}

	return deleted, nil
}


type LogFilters struct {
	UserID        *int64
	ChatID        *int64
	OperationType *model.OperationType
	Level         *model.LogLevel
	Success       *bool
	StartTime     *time.Time
	EndTime       *time.Time
}


type LogStats struct {
	TotalCount    int64
	SuccessCount  int64
	ErrorCount    int64
	AvgDurationMs *float64
}
