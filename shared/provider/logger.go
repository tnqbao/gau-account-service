package provider

import (
	"context"
	"time"

	"github.com/tnqbao/gau-account-service/shared/infra"
)

// LoggerProvider wraps the infra logger to provide a consistent interface
type LoggerProvider struct {
	logger *infra.LoggerClient
}

// NewLoggerProvider creates a new logger provider instance
func NewLoggerProvider() *LoggerProvider {
	return &LoggerProvider{
		logger: infra.GetLogger(),
	}
}

// Context-aware logging methods with trace information
func (lp *LoggerProvider) InfoWithContext(ctx context.Context, msg string, fields map[string]interface{}) {
	lp.logger.InfoWithContext(ctx, msg, fields)
}

func (lp *LoggerProvider) ErrorWithContext(ctx context.Context, msg string, err error, fields map[string]interface{}) {
	lp.logger.ErrorWithContext(ctx, msg, err, fields)
}

func (lp *LoggerProvider) WarningWithContext(ctx context.Context, msg string, fields map[string]interface{}) {
	lp.logger.WarningWithContext(ctx, msg, fields)
}

func (lp *LoggerProvider) DebugWithContext(ctx context.Context, msg string, fields map[string]interface{}) {
	lp.logger.DebugWithContext(ctx, msg, fields)
}

// Printf-style formatting methods with context
func (lp *LoggerProvider) InfoWithContextf(ctx context.Context, format string, args ...interface{}) {
	lp.logger.InfoWithContextf(ctx, format, args...)
}

func (lp *LoggerProvider) ErrorWithContextf(ctx context.Context, err error, format string, args ...interface{}) {
	lp.logger.ErrorWithContextf(ctx, err, format, args...)
}

func (lp *LoggerProvider) WarningWithContextf(ctx context.Context, format string, args ...interface{}) {
	lp.logger.WarningWithContextf(ctx, format, args...)
}

func (lp *LoggerProvider) DebugWithContextf(ctx context.Context, format string, args ...interface{}) {
	lp.logger.DebugWithContextf(ctx, format, args...)
}

// Core logging methods without context
func (lp *LoggerProvider) Info(msg string, fields map[string]interface{}) {
	lp.logger.Info(msg, fields)
}

func (lp *LoggerProvider) Error(msg string, err error, fields map[string]interface{}) {
	lp.logger.Error(msg, err, fields)
}

func (lp *LoggerProvider) Warning(msg string, fields map[string]interface{}) {
	lp.logger.Warning(msg, fields)
}

func (lp *LoggerProvider) Debug(msg string, fields map[string]interface{}) {
	lp.logger.Debug(msg, fields)
}

// Convenience methods for simple logging
func (lp *LoggerProvider) InfoSimple(msg string) {
	lp.logger.InfoSimple(msg)
}

func (lp *LoggerProvider) ErrorSimple(msg string, err error) {
	lp.logger.ErrorSimple(msg, err)
}

func (lp *LoggerProvider) WarningSimple(msg string) {
	lp.logger.WarningSimple(msg)
}

func (lp *LoggerProvider) DebugSimple(msg string) {
	lp.logger.DebugSimple(msg)
}

// HTTP request logging helper
func (lp *LoggerProvider) LogHTTPRequest(ctx context.Context, method, path, userID string, statusCode int, duration time.Duration) {
	lp.logger.LogHTTPRequest(ctx, method, path, userID, statusCode, duration)
}

// Database operation logging helper
func (lp *LoggerProvider) LogDBOperation(ctx context.Context, operation, table string, duration time.Duration, err error) {
	lp.logger.LogDBOperation(ctx, operation, table, duration, err)
}
