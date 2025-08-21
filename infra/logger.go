package infra

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tnqbao/gau-account-service/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LoggerClient struct {
	Logger   *zap.Logger
	Sugar    *zap.SugaredLogger
	shutdown func(context.Context) error
}

var loggerInstance *LoggerClient

func InitLoggerClient(cfg *config.EnvConfig) *LoggerClient {
	if loggerInstance != nil {
		return loggerInstance
	}

	// Initialize OpenTelemetry tracer
	shutdown, err := initTracer(cfg)
	if err != nil {
		log.Printf("Failed to initialize OpenTelemetry tracer: %v", err)
	}

	// Create custom logger configuration
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	zapConfig.OutputPaths = []string{"stdout"}
	zapConfig.ErrorOutputPaths = []string{"stderr"}

	// Add service information to all logs
	zapConfig.InitialFields = map[string]interface{}{
		"service": cfg.Grafana.ServiceName,
		"version": "1.0.0",
	}

	// Customize encoder config for better structured logging
	zapConfig.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	logger, err := zapConfig.Build()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	loggerInstance = &LoggerClient{
		Logger:   logger,
		Sugar:    logger.Sugar(),
		shutdown: shutdown,
	}

	loggerInstance.Info("Logger initialized successfully", map[string]interface{}{
		"grafana_endpoint": cfg.Grafana.OTLPEndpoint,
		"service_name":     cfg.Grafana.ServiceName,
	})

	return loggerInstance
}

func initTracer(cfg *config.EnvConfig) (func(context.Context) error, error) {
	// Create gRPC connection to Grafana OTLP endpoint
	conn, err := grpc.NewClient(
		cfg.Grafana.OTLPEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	// Create OTLP trace exporter
	traceExporter, err := otlptracegrpc.New(
		context.Background(),
		otlptracegrpc.WithGRPCConn(conn),
	)
	if err != nil {
		return nil, err
	}

	// Create resource with service information
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(cfg.Grafana.ServiceName),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		return nil, err
	}

	// Create trace provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}

func GetLogger() *LoggerClient {
	if loggerInstance == nil {
		panic("Logger not initialized. Call InitLoggerClient() first.")
	}
	return loggerInstance
}

func (l *LoggerClient) Shutdown(ctx context.Context) error {
	if err := l.Logger.Sync(); err != nil {
		log.Printf("Failed to sync logger: %v", err)
	}
	if l.shutdown != nil {
		return l.shutdown(ctx)
	}
	return nil
}

// Context-aware logging methods
func (l *LoggerClient) InfoWithContext(ctx context.Context, msg string, fields map[string]interface{}) {
	zapFields := l.extractContextFields(ctx, fields)
	l.Logger.Info(msg, zapFields...)
}

func (l *LoggerClient) ErrorWithContext(ctx context.Context, msg string, err error, fields map[string]interface{}) {
	zapFields := l.extractContextFields(ctx, fields)
	if err != nil {
		zapFields = append(zapFields, zap.Error(err))
	}
	l.Logger.Error(msg, zapFields...)
}

func (l *LoggerClient) WarningWithContext(ctx context.Context, msg string, fields map[string]interface{}) {
	zapFields := l.extractContextFields(ctx, fields)
	l.Logger.Warn(msg, zapFields...)
}

func (l *LoggerClient) DebugWithContext(ctx context.Context, msg string, fields map[string]interface{}) {
	zapFields := l.extractContextFields(ctx, fields)
	l.Logger.Debug(msg, zapFields...)
}

func (l *LoggerClient) FatalWithContext(ctx context.Context, msg string, err error, fields map[string]interface{}) {
	zapFields := l.extractContextFields(ctx, fields)
	if err != nil {
		zapFields = append(zapFields, zap.Error(err))
	}
	l.Logger.Fatal(msg, zapFields...)
}

// printf-style formatting methods
func (l *LoggerClient) InfoWithContextf(ctx context.Context, format string, args ...interface{}) {
	zapFields := l.extractContextFields(ctx, nil)
	msg := fmt.Sprintf(format, args...)
	l.Logger.Info(msg, zapFields...)
}

func (l *LoggerClient) ErrorWithContextf(ctx context.Context, err error, format string, args ...interface{}) {
	zapFields := l.extractContextFields(ctx, nil)
	if err != nil {
		zapFields = append(zapFields, zap.Error(err))
	}
	msg := fmt.Sprintf(format, args...)
	l.Logger.Error(msg, zapFields...)
}

func (l *LoggerClient) WarningWithContextf(ctx context.Context, format string, args ...interface{}) {
	zapFields := l.extractContextFields(ctx, nil)
	msg := fmt.Sprintf(format, args...)
	l.Logger.Warn(msg, zapFields...)
}

func (l *LoggerClient) DebugWithContextf(ctx context.Context, format string, args ...interface{}) {
	zapFields := l.extractContextFields(ctx, nil)
	msg := fmt.Sprintf(format, args...)
	l.Logger.Debug(msg, zapFields...)
}

// Helper method to extract trace information from context
func (l *LoggerClient) extractContextFields(ctx context.Context, fields map[string]interface{}) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields)+3)

	// Add trace information from OpenTelemetry context
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		zapFields = append(zapFields,
			zap.String("trace_id", span.SpanContext().TraceID().String()),
			zap.String("span_id", span.SpanContext().SpanID().String()),
		)
	}

	// Add custom fields
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}

	return zapFields
}

// Core logging methods
func (l *LoggerClient) Info(msg string, fields map[string]interface{}) {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	l.Logger.Info(msg, zapFields...)
}

func (l *LoggerClient) Error(msg string, err error, fields map[string]interface{}) {
	zapFields := make([]zap.Field, 0, len(fields)+1)
	if err != nil {
		zapFields = append(zapFields, zap.Error(err))
	}
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	l.Logger.Error(msg, zapFields...)
}

func (l *LoggerClient) Warning(msg string, fields map[string]interface{}) {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	l.Logger.Warn(msg, zapFields...)
}

func (l *LoggerClient) Debug(msg string, fields map[string]interface{}) {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	l.Logger.Debug(msg, zapFields...)
}

func (l *LoggerClient) Fatal(msg string, err error, fields map[string]interface{}) {
	zapFields := make([]zap.Field, 0, len(fields)+1)
	if err != nil {
		zapFields = append(zapFields, zap.Error(err))
	}
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	l.Logger.Fatal(msg, zapFields...)
}

// Convenience methods for simple logging
func (l *LoggerClient) InfoSimple(msg string) {
	l.Logger.Info(msg)
}

func (l *LoggerClient) ErrorSimple(msg string, err error) {
	if err != nil {
		l.Logger.Error(msg, zap.Error(err))
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
