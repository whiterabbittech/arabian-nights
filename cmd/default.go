package cmd

import (
	"github.com/urfave/cli/v2"
	"github.com/whiterabbittech/arabian-nights/config"
	"github.com/whiterabbittech/arabian-nights/pkg"
)

// Default is the default, top level command for ArabianNights. It is executed
// when the user does not provide a subcommand.
func Default(ctx *cli.Context) error {
	var conf = config.NewConfigFromCLI(ctx)
	var client, err = pkg.NewClient(conf)
	if err != nil {
		return err
	}
	_ = client

	return nil
}
