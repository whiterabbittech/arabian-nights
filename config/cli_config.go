package config

import (
	"time"

	"github.com/urfave/cli/v2"
)

const defaultTimeout = 5 * time.Second

type CLIConfig struct {
	InCluster   bool
	Namespace   string
	ServiceName string
	Timeout     time.Duration
}

func NewConfigFromCLI(ctx *cli.Context) *CLIConfig {
	return &CLIConfig{
		InCluster:   ctx.Bool("in-cluster"),
		Namespace:   ctx.String("namespace"),
		ServiceName: ctx.String("service-name"),
		Timeout:     defaultTimeout,
	}
}
