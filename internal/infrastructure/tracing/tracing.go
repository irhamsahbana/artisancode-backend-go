package tracing

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"strings"

	"codebase-app/internal/infrastructure/config"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/credentials"
)

type Config struct {
	Endpoint   string
	Headers    string
	Insecure   bool
	Debug      bool
	AppName    string
	AppVersion string
	AppEnv     string
	LogWriter  io.Writer
}

var globalServiceName string
var globalLogWriter io.Writer

func InitTracer(cfg *Config) (*sdktrace.TracerProvider, error) {
	globalServiceName = cfg.AppName
	globalLogWriter = cfg.LogWriter

	res, err := resource.New(
		context.Background(),
		resource.WithHost(),
		resource.WithOS(),
		resource.WithAttributes(
			semconv.ServiceName(cfg.AppName),
			semconv.ServiceVersion(cfg.AppVersion),
			semconv.DeploymentEnvironment(cfg.AppEnv),
		),
	)
	if err != nil {
		return nil, err
	}

	tracerProviderOptions := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(res),
	}

	if cfg.Debug {
		exporter, exportErr := stdouttrace.New(stdouttrace.WithPrettyPrint())
		if exportErr != nil {
			return nil, fmt.Errorf("failed to create stdout trace exporter: %w", exportErr)
		}

		tracerProviderOptions = append(tracerProviderOptions, sdktrace.WithBatcher(exporter))
		log.Info().Msg("OpenTelemetry stdout exporter initialized (debug mode)")
	}

	if cfg.Endpoint != "" {
		exporter, exportErr := newExporter(cfg)
		if exportErr != nil {
			return nil, exportErr
		}

		tracerProviderOptions = append(tracerProviderOptions, sdktrace.WithBatcher(exporter))
	} else if !cfg.Debug {
		log.Info().Msg("OpenTelemetry tracer initialized (no exporter, tracing disabled)")
	}

	tp := sdktrace.NewTracerProvider(tracerProviderOptions...)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tp, nil
}

func newExporter(cfg *Config) (sdktrace.SpanExporter, error) {
	headers := parseHeaders(cfg.Headers)

	if strings.HasPrefix(cfg.Endpoint, "http://") || strings.HasPrefix(cfg.Endpoint, "https://") {
		endpoint := strings.TrimPrefix(cfg.Endpoint, "http://")
		endpoint = strings.TrimPrefix(endpoint, "https://")
		endpoint = strings.TrimRight(endpoint, "/")

		opts := []otlptracehttp.Option{
			otlptracehttp.WithEndpoint(endpoint),
		}

		if strings.HasPrefix(cfg.Endpoint, "http://") {
			opts = append(opts, otlptracehttp.WithInsecure())
		}

		if len(headers) > 0 {
			opts = append(opts, otlptracehttp.WithHeaders(headers))
		}

		exporter, err := otlptracehttp.New(context.Background(), opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create otlp http exporter: %w", err)
		}

		log.Info().Str("endpoint", cfg.Endpoint).Msg("OpenTelemetry tracer initialized (OTLP HTTP exporter)")
		return exporter, nil
	}

	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(cfg.Endpoint),
	}

	if cfg.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	} else {
		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
		opts = append(opts, otlptracegrpc.WithTLSCredentials(credentials.NewTLS(tlsConfig)))
	}

	if len(headers) > 0 {
		opts = append(opts, otlptracegrpc.WithHeaders(headers))
	}

	exporter, err := otlptracegrpc.New(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create otlp grpc exporter: %w", err)
	}

	log.Info().Str("endpoint", cfg.Endpoint).Msg("OpenTelemetry tracer initialized (OTLP gRPC exporter)")
	return exporter, nil
}

func parseHeaders(headersStr string) map[string]string {
	headers := make(map[string]string)
	if headersStr == "" {
		return headers
	}

	pairs := strings.Split(headersStr, ",")
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			continue
		}

		headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	return headers
}

func StartSpan(ctx context.Context, name string, opts ...oteltrace.SpanStartOption) (context.Context, oteltrace.Span) {
	tracerName := globalServiceName
	if tracerName == "" {
		tracerName = config.Envs.App.Name
	}

	ctx, span := otel.Tracer(tracerName).Start(ctx, name, opts...)

	if sc := span.SpanContext(); sc.IsValid() {
		logger := log.Logger.With().
			Str("trace_id", sc.TraceID().String()).
			Str("span_id", sc.SpanID().String()).
			Logger()

		ctx = logger.WithContext(ctx)
	}

	return ctx, span
}

func RecordError(span oteltrace.Span, err error) {
	if span == nil || err == nil {
		return
	}

	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

func AddLogEvent(span oteltrace.Span, key string, value string) {
	if span == nil || !span.SpanContext().IsValid() {
		return
	}

	span.AddEvent("log", oteltrace.WithAttributes(attribute.String(key, value)))
}
