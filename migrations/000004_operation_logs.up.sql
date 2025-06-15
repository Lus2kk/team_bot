
CREATE TABLE IF NOT EXISTS operation_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    chat_id BIGINT,
    username VARCHAR(255),
    operation_type VARCHAR(50) NOT NULL,
    level VARCHAR(10) NOT NULL DEFAULT 'INFO',
    message TEXT NOT NULL,
    details TEXT,
    ip_address INET,
    user_agent TEXT,
    success BOOLEAN NOT NULL DEFAULT TRUE,
    duration_ms BIGINT,
    error_code VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);


CREATE INDEX IF NOT EXISTS idx_operation_logs_user_id ON operation_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_operation_logs_chat_id ON operation_logs(chat_id);
CREATE INDEX IF NOT EXISTS idx_operation_logs_operation_type ON operation_logs(operation_type);
CREATE INDEX IF NOT EXISTS idx_operation_logs_level ON operation_logs(level);
CREATE INDEX IF NOT EXISTS idx_operation_logs_created_at ON operation_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_operation_logs_success ON operation_logs(success);


CREATE INDEX IF NOT EXISTS idx_operation_logs_errors ON operation_logs(created_at, operation_type, error_code) 
WHERE success = FALSE;
