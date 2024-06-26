//go:build dev
// +build dev

package main

import (
	_ "github.com/joho/godotenv/autoload"
)

var (
	defaultDataDir              = "data"
	defaultHTTPPort             = 8080
	defaultFetchIntervalSeconds = 5
)
