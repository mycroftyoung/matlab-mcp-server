// Copyright 2026 The MathWorks, Inc.

package provider

import (
	"context"
	"time"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/telemetry/otel"
	"github.com/matlab/matlab-mcp-server/internal/entities"
	"github.com/matlab/matlab-mcp-server/internal/messages"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

const (
	exportTimeout   = 10 * time.Second
	shutdownTimeout = 3 * time.Second
)

type LoggerFactory interface {
	GetGlobalLogger() (entities.Logger, messages.Error)
}

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type LifecycleSignaler interface {
	AddShutdownFunction(shutdownFcn func() error)
}

type Factory struct {
	loggerFactory     LoggerFactory
	configFactory     ConfigFactory
	lifecycleSignaler LifecycleSignaler
}

func NewFactory(
	loggerFactory LoggerFactory,
	configFactory ConfigFactory,
	lifecycleSignaler LifecycleSignaler,
) *Factory {
	return &Factory{
		loggerFactory:     loggerFactory,
		configFactory:     configFactory,
		lifecycleSignaler: lifecycleSignaler,
	}
}

func (f *Factory) New(
	exporter otel.MetricExporter,
	errorHandler otel.ErrorHandler,
	serviceName, serviceVersion string,
) (otel.MeterProvider, messages.Error) {
	logger, messagesErr := f.loggerFactory.GetGlobalLogger()
	if messagesErr != nil {
		return nil, messagesErr
	}

	logger.Debug("Creating OTEL meter provider")
	defer logger.Debug("Done creating OTEL meter provider")

	cfg, messagesErr := f.configFactory.Config()
	if messagesErr != nil {
		return nil, messagesErr
	}

	collectionInterval := cfg.TelemetryCollectionInterval()
	logger.With("interval", collectionInterval.String()).Debug("Configuring telemetry collection interval")

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
		semconv.ServiceVersion(serviceVersion),
	)

	concreteMeterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(
			exporter,
			sdkmetric.WithInterval(collectionInterval),
			sdkmetric.WithTimeout(exportTimeout),
		)),
	)

	f.lifecycleSignaler.AddShutdownFunction(func() error {
		logger.Debug("Shutting down OTEL meter provider")
		defer logger.Debug("Done shutting down OTEL meter provider")

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := concreteMeterProvider.Shutdown(ctx); err != nil {
			errorHandler.Handle(err)
		}
		return nil
	})

	return concreteMeterProvider, nil
}
