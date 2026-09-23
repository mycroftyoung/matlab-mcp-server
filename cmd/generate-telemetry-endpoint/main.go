// Copyright 2026 The MathWorks, Inc.

// generate-telemetry-endpoint writes the build-tagged Go source that bakes the
// telemetry collector endpoint (from the TELEMETRY_COLLECTOR_ENDPOINT
// environment variable) into a release build.
//
// Usage:
//
//	generate-telemetry-endpoint <output-file>
package main

import (
	"fmt"
	"os"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/telemetryendpoint"
)

const endpointEnvVar = "TELEMETRY_COLLECTOR_ENDPOINT"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: generate-telemetry-endpoint <output-file>\n")
		os.Exit(2)
	}

	outputPath := os.Args[1]
	endpoint := os.Getenv(endpointEnvVar)

	if err := telemetryendpoint.Generate(outputPath, endpoint); err != nil {
		fmt.Fprintf(os.Stderr, "generate-telemetry-endpoint: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Wrote telemetry endpoint to %s\n", outputPath)
}
