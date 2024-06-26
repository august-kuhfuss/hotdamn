//go:build tools
// +build tools

package tools

import (
	_ "github.com/bokwoon95/wgo"
	_ "github.com/goreleaser/goreleaser/v2"
	_ "github.com/pressly/goose/v3/cmd/goose"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
