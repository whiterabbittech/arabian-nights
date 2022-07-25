package main

import (
	"log"
	"os"

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
				Name:    "service-name",
				Aliases: []string{"s"},
				Usage:   "the name of the Vault service. Required.",
			},
		},
		Action: func(ctx *cli.Context) error {
			return cmd.Default(ctx)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
