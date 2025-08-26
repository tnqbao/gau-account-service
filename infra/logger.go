package infra

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/tnqbao/gau-account-service/config"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	metricwrap "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
	tracewrap "go.opentelemetry.io/otel/trace"
)

const schemaName = "https://github.com/grafana/docker-otel-lgtm"

// LoggerClient chứa các thành phần cần thiết của OpenTelemetry
type LoggerClient struct {
	Tracer   tracewrap.Tracer
	Logger   *slog.Logger
	Meter    metricwrap.Meter
	shutdown func(context.Context) error
}

var loggerInstance *LoggerClient

func InitLoggerClient(cfg *config.EnvConfig) *LoggerClient {
	if loggerInstance != nil {
		return loggerInstance
	}

	ctx := context.Background()

	// Setup OpenTelemetry SDK
	shutdown, err := setupOTelSDK(ctx, cfg)
	if err != nil {
		slog.Error("Failed to initialize OpenTelemetry SDK", "error", err)
		panic(err)
	}

	// Initialize OpenTelemetry components
	tracer := otel.Tracer(schemaName)
	logger := otelslog.NewLogger(schemaName)
	meter := otel.GetMeterProvider().Meter(schemaName)

	loggerInstance = &LoggerClient{
		Tracer:   tracer,
		Logger:   logger,
		Meter:    meter,
		shutdown: shutdown,
	}

	loggerInstance.Info("Logger initialized successfully", map[string]interface{}{
		"grafana_endpoint": cfg.Grafana.OTLPEndpoint,
		"service_name":     cfg.Grafana.ServiceName,
	})

	return loggerInstance
}

// setupOTelSDK khởi tạo các pipeline cho trace, metric và log của OpenTelemetry
func setupOTelSDK(ctx context.Context, cfg *config.EnvConfig) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown gọi lần lượt các cleanup functions đã đăng ký
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// helper xử lý lỗi, đảm bảo gọi shutdown cho cleanup
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Cấu hình propagator
	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)

	// Cấu hình resource attributes
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.Grafana.ServiceName),
			semconv.DeploymentEnvironmentName(cfg.Environment.Mode), // có thể lấy từ config
			semconv.ServiceNamespace(cfg.Environment.Group),
			attribute.String("service.name", cfg.Grafana.ServiceName),
			attribute.String("deployment.environment", cfg.Environment.Mode),
			attribute.String("service.namespace", cfg.Environment.Group),
		),
	)
	if err != nil {
		handleErr(err)
		return nil, err
	}

	// Thiết lập trace exporter và tracer provider
	traceExporter, err := otlptrace.New(ctx, otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(cfg.Grafana.OTLPEndpoint),
		// Remove WithInsecure() for HTTPS
	))
	if err != nil {
		handleErr(err)
		return nil, err
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
		trace.WithResource(res),
	)

	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Thiết lập metric exporter và meter provider
	metricExporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpoint(cfg.Grafana.OTLPEndpoint),
		// Remove WithInsecure() for HTTPS
	)
	if err != nil {
		handleErr(err)
		return nil, err
	}

	// Set export interval to 5 seconds
	exportInterval := 5 * time.Second
	metricReader := metric.NewPeriodicReader(metricExporter,
		metric.WithInterval(exportInterval),
	)

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metricReader),
		metric.WithResource(res),
	)

	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	// Thiết lập log exporter và logger provider
	logExporter, err := otlploghttp.New(ctx,
		otlploghttp.WithEndpoint(cfg.Grafana.OTLPEndpoint),
		// Remove WithInsecure() for HTTPS
	)
	if err != nil {
		handleErr(err)
		return nil, err
	}

	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
		log.WithResource(res),
	)

	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)

	// Bắt đầu runtime instrumentation
	err = runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
	if err != nil {
		handleErr(err)
		return nil, fmt.Errorf("otel runtime instrumentation failed: %w", err)
	}

	return shutdown, nil
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

// Context-aware logging methods with trace information
func (l *LoggerClient) InfoWithContext(ctx context.Context, msg string, fields map[string]interface{}) {
	attrs := l.extractContextAttributes(ctx, fields)
	// Convert []slog.Attr to []any
	args := make([]any, len(attrs))
	for i, attr := range attrs {
		args[i] = attr
	}
	l.Logger.InfoContext(ctx, msg, args...)
}

func (l *LoggerClient) ErrorWithContext(ctx context.Context, msg string, err error, fields map[string]interface{}) {
	attrs := l.extractContextAttributes(ctx, fields)
	if err != nil {
		attrs = append(attrs, slog.Any("error", err))
	}
	// Convert []slog.Attr to []any
	args := make([]any, len(attrs))
	for i, attr := range attrs {
		args[i] = attr
	}
	l.Logger.ErrorContext(ctx, msg, args...)
}

func (l *LoggerClient) WarningWithContext(ctx context.Context, msg string, fields map[string]interface{}) {
	attrs := l.extractContextAttributes(ctx, fields)
	// Convert []slog.Attr to []any
	args := make([]any, len(attrs))
	for i, attr := range attrs {
		args[i] = attr
	}
	l.Logger.WarnContext(ctx, msg, args...)
}

func (l *LoggerClient) DebugWithContext(ctx context.Context, msg string, fields map[string]interface{}) {
	attrs := l.extractContextAttributes(ctx, fields)
	// Convert []slog.Attr to []any
	args := make([]any, len(attrs))
	for i, attr := range attrs {
		args[i] = attr
	}
	l.Logger.DebugContext(ctx, msg, args...)
}

// Printf-style formatting methods with context
func (l *LoggerClient) InfoWithContextf(ctx context.Context, format string, args ...interface{}) {
	attrs := l.extractContextAttributes(ctx, nil)
	msg := fmt.Sprintf(format, args...)
	// Convert []slog.Attr to []any
	logArgs := make([]any, len(attrs))
	for i, attr := range attrs {
		logArgs[i] = attr
	}
	l.Logger.InfoContext(ctx, msg, logArgs...)
}

func (l *LoggerClient) ErrorWithContextf(ctx context.Context, err error, format string, args ...interface{}) {
	attrs := l.extractContextAttributes(ctx, nil)
	if err != nil {
		attrs = append(attrs, slog.Any("error", err))
	}
	msg := fmt.Sprintf(format, args...)
	// Convert []slog.Attr to []any
	logArgs := make([]any, len(attrs))
	for i, attr := range attrs {
		logArgs[i] = attr
	}
	l.Logger.ErrorContext(ctx, msg, logArgs...)
}

func (l *LoggerClient) WarningWithContextf(ctx context.Context, format string, args ...interface{}) {
	attrs := l.extractContextAttributes(ctx, nil)
	msg := fmt.Sprintf(format, args...)
	// Convert []slog.Attr to []any
	logArgs := make([]any, len(attrs))
	for i, attr := range attrs {
		logArgs[i] = attr
	}
	l.Logger.WarnContext(ctx, msg, logArgs...)
}

func (l *LoggerClient) DebugWithContextf(ctx context.Context, format string, args ...interface{}) {
	attrs := l.extractContextAttributes(ctx, nil)
	msg := fmt.Sprintf(format, args...)
	// Convert []slog.Attr to []any
	logArgs := make([]any, len(attrs))
	for i, attr := range attrs {
		logArgs[i] = attr
	}
	l.Logger.DebugContext(ctx, msg, logArgs...)
}

// Helper method to extract trace information from context
func (l *LoggerClient) extractContextAttributes(ctx context.Context, fields map[string]interface{}) []slog.Attr {
	attrs := make([]slog.Attr, 0, len(fields)+3)

	// Add trace information from OpenTelemetry context
	span := tracewrap.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		attrs = append(attrs,
			slog.String("trace_id", span.SpanContext().TraceID().String()),
			slog.String("span_id", span.SpanContext().SpanID().String()),
		)
	}

	// Add custom fields
	for k, v := range fields {
		attrs = append(attrs, slog.Any(k, v))
	}

	return attrs
}

// Core logging methods without context
func (l *LoggerClient) Info(msg string, fields map[string]interface{}) {
	attrs := make([]any, 0, len(fields))
	for k, v := range fields {
		attrs = append(attrs, slog.Any(k, v))
	}
	l.Logger.Info(msg, attrs...)
}

func (l *LoggerClient) Error(msg string, err error, fields map[string]interface{}) {
	attrs := make([]any, 0, len(fields)+1)
	if err != nil {
		attrs = append(attrs, slog.Any("error", err))
	}
	for k, v := range fields {
		attrs = append(attrs, slog.Any(k, v))
	}
	l.Logger.Error(msg, attrs...)
}

func (l *LoggerClient) Warning(msg string, fields map[string]interface{}) {
	attrs := make([]any, 0, len(fields))
	for k, v := range fields {
		attrs = append(attrs, slog.Any(k, v))
	}
	l.Logger.Warn(msg, attrs...)
}

func (l *LoggerClient) Debug(msg string, fields map[string]interface{}) {
	attrs := make([]any, 0, len(fields))
	for k, v := range fields {
		attrs = append(attrs, slog.Any(k, v))
	}
	l.Logger.Debug(msg, attrs...)
}

// Convenience methods for simple logging
func (l *LoggerClient) InfoSimple(msg string) {
	l.Logger.Info(msg)
}

func (l *LoggerClient) ErrorSimple(msg string, err error) {
	if err != nil {
		l.Logger.Error(msg, slog.Any("error", err))
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
func (l *LoggerClient) LogHTTPRequest(ctx context.Context, method, path, userID string, statusCode int, duration time.Duration) {
	l.InfoWithContext(ctx, "HTTP Request", map[string]interface{}{
		"method":      method,
		"path":        path,
		"user_id":     userID,
		"status_code": statusCode,
		"duration_ms": duration.Milliseconds(),
	})
}

// Database operation logging helper
func (l *LoggerClient) LogDBOperation(ctx context.Context, operation, table string, duration time.Duration, err error) {
	fields := map[string]interface{}{
		"operation":   operation,
		"table":       table,
		"duration_ms": duration.Milliseconds(),
	}

	if err != nil {
		l.ErrorWithContext(ctx, "Database operation failed", err, fields)
	} else {
		l.InfoWithContext(ctx, "Database operation completed", fields)
	}
}
