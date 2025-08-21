package provider

import (
	"context"
	"time"

	"github.com/tnqbao/gau-account-service/infra"
)

type LoggerProvider struct {
	logger *infra.LoggerClient
}

func NewLoggerProvider() *LoggerProvider {
	return &LoggerProvider{
		logger: infra.GetLogger(),
	}
}

// Info logs an informational message with optional fields
func (lp *LoggerProvider) Info(msg string, fields ...map[string]interface{}) {
	var mergedFields map[string]interface{}
	if len(fields) > 0 {
		mergedFields = fields[0]
	} else {
		mergedFields = make(map[string]interface{})
	}
	lp.logger.Info(msg, mergedFields)
}

// Error logs an error message with optional fields
func (lp *LoggerProvider) Error(msg string, err error, fields ...map[string]interface{}) {
	var mergedFields map[string]interface{}
	if len(fields) > 0 {
		mergedFields = fields[0]
	} else {
		mergedFields = make(map[string]interface{})
	}
	lp.logger.Error(msg, err, mergedFields)
}

// Warning logs a warning message with optional fields
func (lp *LoggerProvider) Warning(msg string, fields ...map[string]interface{}) {
	var mergedFields map[string]interface{}
	if len(fields) > 0 {
		mergedFields = fields[0]
	} else {
		mergedFields = make(map[string]interface{})
	}
	lp.logger.Warning(msg, mergedFields)
}

// Debug logs a debug message with optional fields
func (lp *LoggerProvider) Debug(msg string, fields ...map[string]interface{}) {
	var mergedFields map[string]interface{}
	if len(fields) > 0 {
		mergedFields = fields[0]
	} else {
		mergedFields = make(map[string]interface{})
	}
	lp.logger.Debug(msg, mergedFields)
}

// Fatal logs a fatal message and exits the application
func (lp *LoggerProvider) Fatal(msg string, err error, fields ...map[string]interface{}) {
	var mergedFields map[string]interface{}
	if len(fields) > 0 {
		mergedFields = fields[0]
	} else {
		mergedFields = make(map[string]interface{})
	}
	lp.logger.Fatal(msg, err, mergedFields)
}

// Simple logging methods without fields
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

// Specialized logging methods for common use cases

// LogUserAction logs user-related actions
func (lp *LoggerProvider) LogUserAction(userID, action, details string) {
	lp.Info("User action", map[string]interface{}{
		"user_id": userID,
		"action":  action,
		"details": details,
	})
}

// LogAuthEvent logs authentication-related events
func (lp *LoggerProvider) LogAuthEvent(userID, event, method string, success bool) {
	lp.Info("Authentication event", map[string]interface{}{
		"user_id": userID,
		"event":   event,
		"method":  method,
		"success": success,
	})
}

// LogMFAEvent logs MFA-related events
func (lp *LoggerProvider) LogMFAEvent(userID, event string, success bool, details string) {
	lp.Info("MFA event", map[string]interface{}{
		"user_id": userID,
		"event":   event,
		"success": success,
		"details": details,
	})
}

// LogAPIRequest logs API request details
func (lp *LoggerProvider) LogAPIRequest(method, path, userID string, statusCode int, duration time.Duration) {
	lp.logger.LogHTTPRequest(method, path, userID, statusCode, duration)
}

// LogDatabaseOperation logs database operation details
func (lp *LoggerProvider) LogDatabaseOperation(operation, table string, duration time.Duration, err error) {
	lp.logger.LogDBOperation(operation, table, duration, err)
}

// LogSecurityEvent logs security-related events
func (lp *LoggerProvider) LogSecurityEvent(userID, event, details string, severity string) {
	fields := map[string]interface{}{
		"user_id":  userID,
		"event":    event,
		"details":  details,
		"severity": severity,
		"category": "security",
	}

	switch severity {
	case "high", "critical":
		lp.Error("Security event", nil, fields)
	case "medium":
		lp.Warning("Security event", fields)
	default:
		lp.Info("Security event", fields)
	}
}

// LogProfileUpdate logs user profile update events
func (lp *LoggerProvider) LogProfileUpdate(userID, field, oldValue, newValue string) {
	lp.Info("Profile updated", map[string]interface{}{
		"user_id":   userID,
		"field":     field,
		"old_value": oldValue,
		"new_value": newValue,
	})
}

// LogValidationError logs validation errors
func (lp *LoggerProvider) LogValidationError(userID, field, value, reason string) {
	lp.Warning("Validation error", map[string]interface{}{
		"user_id": userID,
		"field":   field,
		"value":   value,
		"reason":  reason,
	})
}

// LogServiceCall logs external service calls
func (lp *LoggerProvider) LogServiceCall(service, endpoint string, duration time.Duration, statusCode int, err error) {
	fields := map[string]interface{}{
		"service":     service,
		"endpoint":    endpoint,
		"duration_ms": duration.Milliseconds(),
		"status_code": statusCode,
	}

	if err != nil {
		lp.Error("Service call failed", err, fields)
	} else {
		lp.Info("Service call completed", fields)
	}
}

// LogBusinessEvent logs business logic events
func (lp *LoggerProvider) LogBusinessEvent(event, details string, metadata map[string]interface{}) {
	fields := map[string]interface{}{
		"event":    event,
		"details":  details,
		"category": "business",
	}

	// Merge metadata
	for k, v := range metadata {
		fields[k] = v
	}

	lp.Info("Business event", fields)
}

// Context-aware logging methods
func (lp *LoggerProvider) InfoWithContext(ctx context.Context, msg string, fields ...map[string]interface{}) {
	var mergedFields map[string]interface{}
	if len(fields) > 0 {
		mergedFields = fields[0]
	} else {
		mergedFields = make(map[string]interface{})
	}
	lp.logger.InfoWithContext(ctx, msg, mergedFields)
}

func (lp *LoggerProvider) ErrorWithContext(ctx context.Context, msg string, err error, fields ...map[string]interface{}) {
	var mergedFields map[string]interface{}
	if len(fields) > 0 {
		mergedFields = fields[0]
	} else {
		mergedFields = make(map[string]interface{})
	}
	lp.logger.ErrorWithContext(ctx, msg, err, mergedFields)
}

func (lp *LoggerProvider) WarningWithContext(ctx context.Context, msg string, fields ...map[string]interface{}) {
	var mergedFields map[string]interface{}
	if len(fields) > 0 {
		mergedFields = fields[0]
	} else {
		mergedFields = make(map[string]interface{})
	}
	lp.logger.WarningWithContext(ctx, msg, mergedFields)
}

func (lp *LoggerProvider) DebugWithContext(ctx context.Context, msg string, fields ...map[string]interface{}) {
	var mergedFields map[string]interface{}
	if len(fields) > 0 {
		mergedFields = fields[0]
	} else {
		mergedFields = make(map[string]interface{})
	}
	lp.logger.DebugWithContext(ctx, msg, mergedFields)
}

func (lp *LoggerProvider) FatalWithContext(ctx context.Context, msg string, err error, fields ...map[string]interface{}) {
	var mergedFields map[string]interface{}
	if len(fields) > 0 {
		mergedFields = fields[0]
	} else {
		mergedFields = make(map[string]interface{})
	}
	lp.logger.FatalWithContext(ctx, msg, err, mergedFields)
}

// InfoWithContextf Context-aware logging methods with printf-style formatting
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

// Shutdown gracefully shuts down the logger
func (lp *LoggerProvider) Shutdown(ctx context.Context) error {
	return lp.logger.Shutdown(ctx)
}
