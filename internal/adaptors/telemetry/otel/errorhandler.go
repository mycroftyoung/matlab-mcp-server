// Copyright 2026 The MathWorks, Inc.

package otel

import (
	"sync"

	"github.com/matlab/matlab-mcp-server/internal/entities"
	otelapi "go.opentelemetry.io/otel"
)

func NewErrorHandler(logger entities.Logger) ErrorHandler {
	return &loggerErrorHandler{logger: logger}
}

func SetErrorHandler(handler ErrorHandler) {
	otelapi.SetErrorHandler(handler)
}

type loggerErrorHandler struct {
	logger   entities.Logger
	warnOnce sync.Once
}

// An unreachable collector errors on every export, so warn only once.
func (h *loggerErrorHandler) Handle(err error) {
	h.warnOnce.Do(func() {
		h.logger.Warn("Telemetry error; metrics may not be reaching the collector. Further telemetry errors are logged at debug level.")
	})
	h.logger.WithError(err).Debug("Telemetry error")
}
