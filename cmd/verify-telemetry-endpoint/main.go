// Copyright 2026 The MathWorks, Inc.

// verify-telemetry-endpoint fails a release build whose binary is missing the
// TELEMETRY_COLLECTOR_ENDPOINT it was asked to bake in.
//
// Usage:
//
//	verify-telemetry-endpoint <binary-file>
package main

import (
	"fmt"
	"os"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/telemetryendpoint"
)

const endpointEnvVar = "TELEMETRY_COLLECTOR_ENDPOINT"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: verify-telemetry-endpoint <binary-file>\n")
		os.Exit(2)
	}

	binaryPath := os.Args[1]
	endpoint := os.Getenv(endpointEnvVar)

	if err := telemetryendpoint.Verify(binaryPath, endpoint); err != nil {
		fmt.Fprintf(os.Stderr, "verify-telemetry-endpoint: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Verified telemetry endpoint in %s\n", binaryPath)
}
