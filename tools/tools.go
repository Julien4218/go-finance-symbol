//go:build tools
// +build tools

package main

import (
	// build/test.mk
	_ "github.com/stretchr/testify/assert"

	_ "github.com/psampaz/go-mod-outdated"
	_ "golang.org/x/tools/cmd/goimports"

	// build/test.mk
	_ "gotest.tools/gotestsum"
)
