//go:build dev
// +build dev

package main

import "github.com/urfave/cli/v2"

func init() {
	cmd.Flags = append(cmd.Flags, &cli.BoolFlag{
		Usage:   "enable demo mode (mocked data, no tasks)",
		Name:    "enable-demo-mode",
		Aliases: []string{"demo"},
		EnvVars: []string{"ENABLE_DEMO"},
		Value:   false,
		Action: func(ctx *cli.Context, b bool) error {
			demoMode = b
			return nil
		},
	})
}
