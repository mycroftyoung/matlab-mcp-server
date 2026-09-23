// Copyright 2026 The MathWorks, Inc.

package otel_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/telemetry/otel"
	"github.com/matlab/matlab-mcp-server/internal/testutils"
	entitiesmocks "github.com/matlab/matlab-mcp-server/mocks/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestErrorHandler_Handle_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	handler := otel.NewErrorHandler(mockLogger)

	// Act
	handler.Handle(assert.AnError)

	// Assert
	warnLogs := mockLogger.WarnLogs()
	require.Len(t, warnLogs, 1)
	for _, fields := range warnLogs {
		_, hasError := fields["error"]
		assert.False(t, hasError, "warning must not carry the error, which can contain the endpoint URL")
	}

	debugLogs := mockLogger.DebugLogs()
	require.Len(t, debugLogs, 1)
	for _, fields := range debugLogs {
		assert.Equal(t, assert.AnError, fields["error"])
	}
}

func TestErrorHandler_Handle_DoesNotLogEndpointInWarning(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	handler := otel.NewErrorHandler(mockLogger)

	const endpoint = "https://telemetry.example.invalid/v1/metrics"
	exportErr := errors.New(`export failed: Post "` + endpoint + `": dial tcp: connection refused`)

	// Act
	handler.Handle(exportErr)

	// Assert
	warnLogs := mockLogger.WarnLogs()
	require.Len(t, warnLogs, 1)
	for msg, fields := range warnLogs {
		assert.NotContains(t, msg, endpoint)
		for _, value := range fields {
			assert.NotContains(t, fmt.Sprintf("%v", value), endpoint)
		}
	}
}

func TestErrorHandler_Handle_SurfacesWarningOnce(t *testing.T) {
	// Arrange
	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	firstErr := errors.New("first export failed")
	secondErr := errors.New("second export failed")

	mockLogger.EXPECT().Warn(mock.Anything).Once()
	mockLogger.EXPECT().WithError(firstErr).Return(mockLogger).Once()
	mockLogger.EXPECT().WithError(secondErr).Return(mockLogger).Once()
	mockLogger.EXPECT().Debug(mock.Anything).Times(2)

	handler := otel.NewErrorHandler(mockLogger)

	// Act
	handler.Handle(firstErr)
	handler.Handle(secondErr)
}
