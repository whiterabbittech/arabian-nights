package main

import (
	"os"
	"sort"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/whiterabbittech/arabian-nights/cmd"
)

// TODO:
// Determine if we're running in cluster or out of cluster.
// This changes how we authenticate.
// If out-of-cluster, grab the active Kubeconfig.
// If in-cluster, grab a ServiceAccount token.

// TODO: Respect CLI flag to take the path to Kubeconfig file.
//       Respect $KUBECONFIG env variable
//       Have the same behaivior as `kubectl config --help` describes.
func main() {
	app := &cli.App{
		Name:  "arabian-nights",
		Usage: "Unseal Vault and store the unseal keys in a secret",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "color",
				Aliases: []string{"c"},
				Value:   "line",
				Usage:   "determines the behavior of log coloring. \"line\" will color the whole line, and is best for interactive use. \"auto\" detects if stdout is a tty and behaves like \"on\" if so, and \"off\" if not. \"off\" disables colors, and \"on\" only colors the log level.",
			},
			&cli.BoolFlag{
				Name:    "in-cluster",
				Value:   true,
				Aliases: []string{"i"},
				Usage:   "if true, authentication is done with a ServiceAccount token. Otherwise, authentication is done with local kubeconfig",
			},
			&cli.StringFlag{
				Name:    "namespace",
				Aliases: []string{"n"},
				Value:   "",
				Usage:   "the namespace of the Vault Service",
			},
			&cli.StringFlag{
				Name:    "logger-type",
				Aliases: []string{"lt"},
				Value:   "text",
				Usage:   "determines the format of log lines",
			},
			&cli.StringFlag{
				Name:    "log-level",
				Aliases: []string{"v"},
				Value:   "info",
				Usage:   "determines the log level emitted. Emits logs at or below the provided level. Permitted values in ascending order are \"trace\", \"debug\", \"info\", \"warn\", \"error\", and \"fatal\"",
			},
			// TODO: Accept this as an argument instead of a flag.
			&cli.StringFlag{
				Name:     "service-name",
				Aliases:  []string{"s"},
				Usage:    "the name of the Vault service.",
				Required: true,
			},
		},
		Action: func(ctx *cli.Context) error {
			return cmd.Default(ctx)
		},
	}
	sort.Sort(cli.FlagsByName(app.Flags))

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
