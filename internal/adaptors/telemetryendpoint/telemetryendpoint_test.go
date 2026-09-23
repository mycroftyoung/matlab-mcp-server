// Copyright 2026 The MathWorks, Inc.

package telemetryendpoint_test

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/telemetryendpoint"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate_HappyPath(t *testing.T) {
	// Arrange
	outputPath := filepath.Join(t.TempDir(), "endpoint_generated.go")
	const endpoint = "https://telemetry.example.test/v1/metrics"

	// Act
	err := telemetryendpoint.Generate(outputPath, endpoint)

	// Assert
	require.NoError(t, err)

	source, err := os.ReadFile(outputPath) //nolint:gosec // test controls the path
	require.NoError(t, err)

	content := string(source)
	assert.Contains(t, content, "//go:build release")
	assert.Contains(t, content, "package defaultparameters")
	assert.Contains(t, content, "const defaultTelemetryCollectorEndpoint = "+strconv.Quote(endpoint))
}

func TestGenerate_EmptyEndpointReturnsError(t *testing.T) {
	// Arrange
	outputPath := filepath.Join(t.TempDir(), "endpoint_generated.go")

	// Act
	err := telemetryendpoint.Generate(outputPath, "")

	// Assert
	require.Error(t, err)
	assert.NoFileExists(t, outputPath)
}

func TestGenerate_EndpointWithoutSchemeReturnsError(t *testing.T) {
	// Arrange
	outputPath := filepath.Join(t.TempDir(), "endpoint_generated.go")

	// Act
	err := telemetryendpoint.Generate(outputPath, "telemetry.example.test/v1/metrics")

	// Assert
	require.Error(t, err)
	assert.NoFileExists(t, outputPath)
}

func TestGenerate_NonHTTPSchemeReturnsError(t *testing.T) {
	// Arrange
	outputPath := filepath.Join(t.TempDir(), "endpoint_generated.go")

	// Act
	err := telemetryendpoint.Generate(outputPath, "ftp://telemetry.example.test/v1/metrics")

	// Assert
	require.Error(t, err)
	assert.NoFileExists(t, outputPath)
}

func TestVerify_HappyPath(t *testing.T) {
	// Arrange
	const endpoint = "https://telemetry.example.test/v1/metrics"
	binaryPath := filepath.Join(t.TempDir(), "matlab-mcp-server")
	require.NoError(t, os.WriteFile(binaryPath, []byte("padding "+endpoint+" padding"), 0o600))

	// Act
	err := telemetryendpoint.Verify(binaryPath, endpoint)

	// Assert
	require.NoError(t, err)
}

func TestVerify_BinaryWithoutEndpointReturnsError(t *testing.T) {
	// Arrange
	binaryPath := filepath.Join(t.TempDir(), "matlab-mcp-server")
	require.NoError(t, os.WriteFile(binaryPath, []byte("padding"), 0o600))

	// Act
	err := telemetryendpoint.Verify(binaryPath, "https://telemetry.example.test/v1/metrics")

	// Assert
	require.Error(t, err)
}

func TestVerify_EmptyEndpointReturnsErrorInsteadOfMatchingAnyBinary(t *testing.T) {
	// Arrange
	binaryPath := filepath.Join(t.TempDir(), "matlab-mcp-server")
	require.NoError(t, os.WriteFile(binaryPath, []byte("padding"), 0o600))

	// Act
	err := telemetryendpoint.Verify(binaryPath, "")

	// Assert
	require.Error(t, err)
}

func TestVerify_MissingBinaryReturnsError(t *testing.T) {
	// Arrange
	binaryPath := filepath.Join(t.TempDir(), "matlab-mcp-server")

	// Act
	err := telemetryendpoint.Verify(binaryPath, "https://telemetry.example.test/v1/metrics")

	// Assert
	require.Error(t, err)
}
