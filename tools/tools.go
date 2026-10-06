//go:build generate

// Package tools tracks generation and linting tools.
package tools

import (
	_ "github.com/golangci/golangci-lint/v2/cmd/golangci-lint"
	_ "github.com/hashicorp/copywrite"
)

// Generate copyright headers
//go:generate go run github.com/hashicorp/copywrite headers -d .. --config ../.copywrite.hcl

// Run golangci-lint
//go:generate sh -c "cd .. && golangci-lint run ./..."
