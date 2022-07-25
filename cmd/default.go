package cmd

import (
	"context"

	"github.com/urfave/cli/v2"
	"github.com/whiterabbittech/arabian-nights/config"
	"github.com/whiterabbittech/arabian-nights/pkg"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Default is the default, top level command for ArabianNights. It is executed
// when the user does not provide a subcommand.
func Default(cliCtx *cli.Context) error {
	var conf = config.NewConfigFromCLI(cliCtx)
	var ctx = context.Background()

	var client, err = pkg.NewClient(conf)
	if err != nil {
		return err
	}
	// Fetch the Vault Service.
	service, err := getVaultService(ctx, conf, client)
	if err != nil {
		return err
	}
	_ = service

	return nil
}

func getVaultService(ctx context.Context, conf *config.CLIConfig, client *kubernetes.Clientset) (*v1.Service, error) {
	ctx, cancel := context.WithTimeout(ctx, conf.Timeout)
	defer cancel()
	return client.
		CoreV1().
		Services(conf.Namespace).
		Get(ctx, conf.ServiceName, metav1.GetOptions{})
}
