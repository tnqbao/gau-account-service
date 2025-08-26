package infra

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/tnqbao/gau-account-service/config"
)

const schemaName = "gau-account-service"

// LoggerClient tương thích với code cũ
type LoggerClient struct {
	Logger   *slog.Logger
	shutdown func(context.Context) error
}

var loggerInstance *LoggerClient

// SetupOTelLogSDK khởi tạo structured logger với slog
func SetupOTelLogSDK(ctx context.Context, cfg *config.EnvConfig) (shutdown func(context.Context) error, err error) {
	// Tạo structured logger với slog
	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}

	// Tạo JSON handler để xuất structured logs
	handler := slog.NewJSONHandler(os.Stdout, opts)

	// Wrap với context để thêm service metadata
	logger := slog.New(handler).With(
		"service", cfg.Grafana.ServiceName,
		"schema", schemaName,
	)

	// Đặt làm default logger
	slog.SetDefault(logger)

	// Return empty shutdown function
	shutdown = func(ctx context.Context) error {
		return nil
	}

	return shutdown, nil
}

// InitLoggerClient khởi tạo logger client (tương thích với code cũ)
func InitLoggerClient(cfg *config.EnvConfig) *LoggerClient {
	if loggerInstance != nil {
		return loggerInstance
	}

	// Initialize structured logging
	shutdown, err := SetupOTelLogSDK(context.Background(), cfg)
	if err != nil {
		fmt.Printf("Failed to initialize structured logging: %v", err)
		// Fallback to basic slog if setup fails
		logger := slog.Default()
		loggerInstance = &LoggerClient{
			Logger:   logger,
			shutdown: nil,
		}
		return loggerInstance
	}

	// Create structured logger với service metadata
	logger := slog.Default().With(
		"service", cfg.Grafana.ServiceName,
		"endpoint", cfg.Grafana.OTLPEndpoint,
	)

	loggerInstance = &LoggerClient{
		Logger:   logger,
		shutdown: shutdown,
	}

	loggerInstance.Info("Logger initialized successfully", map[string]interface{}{
		"grafana_endpoint": cfg.Grafana.OTLPEndpoint,
		"service_name":     cfg.Grafana.ServiceName,
	})

	return loggerInstance
}

func GetLogger() *LoggerClient {
	if loggerInstance == nil {
		panic("Logger not initialized. Call InitLoggerClient() first.")
	}
	return loggerInstance
}

func (l *LoggerClient) Shutdown(ctx context.Context) error {
	if l.shutdown != nil {
		return l.shutdown(ctx)
	}
	return nil
}

// Context-aware logging methods
func (l *LoggerClient) InfoWithContext(ctx context.Context, msg string, fields map[string]interface{}) {
	attrs := l.convertToSlogAny(fields)
	l.Logger.InfoContext(ctx, msg, attrs...)
}

func (l *LoggerClient) ErrorWithContext(ctx context.Context, msg string, err error, fields map[string]interface{}) {
	attrs := l.convertToSlogAny(fields)
	if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))
	}
	l.Logger.ErrorContext(ctx, msg, attrs...)
}

func (l *LoggerClient) WarningWithContext(ctx context.Context, msg string, fields map[string]interface{}) {
	attrs := l.convertToSlogAny(fields)
	l.Logger.WarnContext(ctx, msg, attrs...)
}

func (l *LoggerClient) DebugWithContext(ctx context.Context, msg string, fields map[string]interface{}) {
	attrs := l.convertToSlogAny(fields)
	l.Logger.DebugContext(ctx, msg, attrs...)
}

func (l *LoggerClient) FatalWithContext(ctx context.Context, msg string, err error, fields map[string]interface{}) {
	attrs := l.convertToSlogAny(fields)
	if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))
	}
	l.Logger.ErrorContext(ctx, msg, attrs...)
	panic(msg) // Fatal should terminate the program
}

// printf-style formatting methods
func (l *LoggerClient) InfoWithContextf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.Logger.InfoContext(ctx, msg)
}

func (l *LoggerClient) ErrorWithContextf(ctx context.Context, err error, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if err != nil {
		l.Logger.ErrorContext(ctx, msg, slog.String("error", err.Error()))
	} else {
		l.Logger.ErrorContext(ctx, msg)
	}
}

func (l *LoggerClient) WarningWithContextf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.Logger.WarnContext(ctx, msg)
}

func (l *LoggerClient) DebugWithContextf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.Logger.DebugContext(ctx, msg)
}

// Helper method to convert map to slog any attributes (compatible with slog API)
func (l *LoggerClient) convertToSlogAny(fields map[string]interface{}) []any {
	attrs := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		attrs = append(attrs, k, v)
	}
	return attrs
}

// Core logging methods
func (l *LoggerClient) Info(msg string, fields map[string]interface{}) {
	attrs := l.convertToSlogAny(fields)
	l.Logger.Info(msg, attrs...)
}

func (l *LoggerClient) Error(msg string, err error, fields map[string]interface{}) {
	attrs := l.convertToSlogAny(fields)
	if err != nil {
		attrs = append(attrs, "error", err.Error())
	}
	l.Logger.Error(msg, attrs...)
}

func (l *LoggerClient) Warning(msg string, fields map[string]interface{}) {
	attrs := l.convertToSlogAny(fields)
	l.Logger.Warn(msg, attrs...)
}

func (l *LoggerClient) Debug(msg string, fields map[string]interface{}) {
	attrs := l.convertToSlogAny(fields)
	l.Logger.Debug(msg, attrs...)
}

func (l *LoggerClient) Fatal(msg string, err error, fields map[string]interface{}) {
	attrs := l.convertToSlogAny(fields)
	if err != nil {
		attrs = append(attrs, "error", err.Error())
	}
	l.Logger.Error(msg, attrs...)
	panic(msg) // Fatal should terminate the program
}

// Convenience methods for simple logging
func (l *LoggerClient) InfoSimple(msg string) {
	l.Logger.Info(msg)
}

func (l *LoggerClient) ErrorSimple(msg string, err error) {
	if err != nil {
		l.Logger.Error(msg, "error", err.Error())
	} else {
		l.Logger.Error(msg)
	}
}

func (l *LoggerClient) WarningSimple(msg string) {
	l.Logger.Warn(msg)
}

func (l *LoggerClient) DebugSimple(msg string) {
	l.Logger.Debug(msg)
}

// HTTP request logging helper
func (l *LoggerClient) LogHTTPRequest(method, path, userID string, statusCode int, duration time.Duration) {
	l.Info("HTTP Request", map[string]interface{}{
		"method":      method,
		"path":        path,
		"user_id":     userID,
		"status_code": statusCode,
		"duration_ms": duration.Milliseconds(),
	})
}

// Database operation logging helper
func (l *LoggerClient) LogDBOperation(operation, table string, duration time.Duration, err error) {
	fields := map[string]interface{}{
		"operation":   operation,
		"table":       table,
		"duration_ms": duration.Milliseconds(),
	}

	if err != nil {
		l.Error("Database operation failed", err, fields)
	} else {
		l.Info("Database operation completed", fields)
	}
}
