package config

import (
	"github.com/urfave/cli/v2"
)

type CLIConfig struct {
	InCluster bool
}

func NewConfigFromCLI(ctx *cli.Context) *CLIConfig {
	return &CLIConfig{
		InCluster: ctx.Bool("in-cluster"),
	}
}
