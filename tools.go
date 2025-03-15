//go:build tools
// +build tools

package tool

import (
	_ "github.com/golang/mock/mockgen"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/swaggo/swag/cmd/swag"
)
