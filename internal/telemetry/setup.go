package telemetry

import (
	"context"
	"errors"
	"time"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
)

func SetupOTelSDK(ctx context.Context, serviceName string, envConfig *config.EnvConfig) (shutdown func(context.Context) error, tp trace.TracerProvider, mp metric.MeterProvider, propogator propagation.TextMapPropagator, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := NewPropagator()
	otel.SetTextMapPropagator(prop)

	res, err := NewResource(ctx)
	if err != nil {
		return shutdown, nil, nil, nil, err
	}
	sres, err := resource.Merge(res, resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
	))

	// Set up trace provider.
	tracerProvider, tracerProviderShutdown, err := NewTraceProvider(ctx, sres, envConfig)
	if err != nil {
		handleErr(err)
		return
	}
	if tracerProviderShutdown != nil {
		shutdownFuncs = append(shutdownFuncs, tracerProviderShutdown)
	}
	otel.SetTracerProvider(tracerProvider)

	// Set up meter provider.
	meterProvider, meterProviderShutdown, err := NewMeterProvider(ctx, sres, envConfig)
	if err != nil {
		handleErr(err)
		return
	}
	if meterProviderShutdown != nil {
		shutdownFuncs = append(shutdownFuncs, meterProviderShutdown)
	}
	otel.SetMeterProvider(meterProvider)

	return shutdown, tracerProvider, meterProvider, prop, nil
}

func NewResource(ctx context.Context) (*resource.Resource, error) {
	res, err := resource.New(ctx,
		resource.WithOS(),
		resource.WithHost(),
		resource.WithProcess(),
	)
	return res, err
}

func NewPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func NewTraceProvider(ctx context.Context, res *resource.Resource, envConfig *config.EnvConfig) (trace.TracerProvider, func(context.Context) error, error) {
	if !envConfig.ENABLE_TELEMETRY {
		return tracenoop.NewTracerProvider(), nil, nil
	}
	traceExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(envConfig.OTLP_ENDPOINT),
	)
	if err != nil {
		return nil, nil, err
	}

	traceProvider := tracesdk.NewTracerProvider(
		tracesdk.WithResource(res),
		tracesdk.WithBatcher(traceExporter,
			tracesdk.WithBatchTimeout(5*time.Second)),
	)
	return traceProvider, traceProvider.Shutdown, nil
}

func NewMeterProvider(ctx context.Context, res *resource.Resource, envConfig *config.EnvConfig) (metric.MeterProvider, func(context.Context) error, error) {
	if envConfig.ENABLE_TELEMETRY {
		metricExporter, err := otlpmetricgrpc.New(
			ctx,
			otlpmetricgrpc.WithInsecure(),
			otlpmetricgrpc.WithEndpoint(envConfig.OTLP_ENDPOINT),
		)
		if err != nil {
			return nil, nil, err
		}

		meterProvider := metricsdk.NewMeterProvider(
			metricsdk.WithResource(res),
			metricsdk.WithReader(metricsdk.NewPeriodicReader(metricExporter,
				metricsdk.WithInterval(1*time.Minute))),
		)
		return meterProvider, meterProvider.Shutdown, nil
	} else {
		return metricnoop.NewMeterProvider(), nil, nil
	}
}
