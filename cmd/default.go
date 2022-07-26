package cmd

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/whiterabbittech/arabian-nights/config"
	"github.com/whiterabbittech/arabian-nights/pkg"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Default is the default, top level command for ArabianNights. It is executed
// when the user does not provide a subcommand.
type DefaultCmd struct {
	ctx  context.Context
	conf *config.CLIConfig
}

func NewDefaultCmd(cliCtx *cli.Context) (*DefaultCmd, error) {
	var conf, err = config.NewConfigFromCLI(cliCtx)
	if err != nil {
		return nil, err
	}
	var cmd = &DefaultCmd{
		ctx:  context.Background(),
		conf: conf,
	}

	return cmd, nil
}

func (cmd *DefaultCmd) Ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(cmd.ctx, cmd.conf.Timeout)
}

func (cmd *DefaultCmd) Run() error {
	client, err := pkg.NewClient(cmd.conf)
	if err != nil {
		return err
	}
	// Fetch the Vault Service.
	service, err := cmd.getVaultService(client)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Expected to find service %s, but no service was found.", cmd.conf.ServiceName)
	}

	var serviceSelector = service.Spec.Selector
	var podSelectorQuery = joinSelectorIntoQuery(serviceSelector)
	_ = podSelectorQuery
	// This link shows to how to get the pods for this service.
	// https://stackoverflow.com/questions/60221012/how-to-list-names-of-all-pods-serving-traffic-behind-a-service-in-kubernetes

	// TODO: this is where I've left off.
	// https://pkg.go.dev/k8s.io/api/core/v1#PodList

	return nil

}

func (cmd *DefaultCmd) getVaultService(client *kubernetes.Clientset) (*v1.Service, error) {
	ctx, cancel := cmd.Ctx()
	defer cancel()
	return client.
		CoreV1().
		Services(cmd.conf.Namespace).
		Get(ctx, cmd.conf.ServiceName, metav1.GetOptions{})
}

// This converts a select of type map[string]string into a selector of type string
// by joining each k/v with a comma, which is the format accepted by the CLI.
func joinSelectorIntoQuery(inputSelector map[string]string) string {
	var keys = getSortedKeys(inputSelector)
	// Map each key to a value.
	var boundKVs = bindKV(keys, inputSelector)
	// Finally, we can join these pairs with commas.
	return strings.Join(boundKVs, ",")
}

func getSortedKeys(selector map[string]string) []string {
	var keys = make([]string, 0, len(selector))
	for key := range selector {
		keys = append(keys, key)
	}
	// Sort the keys alphabetically.
	sort.Strings(keys)
	return keys
}

// Binds each key in keys to its value in selector using '='
func bindKV(keys []string, selector map[string]string) []string {
	var pairs = make([]string, 0, len(keys))
	for _, key := range keys {
		var sb strings.Builder
		var val = selector[key]
		sb.WriteString(key)
		sb.WriteRune('=')
		sb.WriteString(val)
		pairs = append(pairs, sb.String())
	}
	return pairs
}

func (cmd *DefaultCmd) getVaultPods(client *kubernetes.Clientset, selector string) (*v1.PodList, error) {
	ctx, cancel := cmd.Ctx()
	defer cancel()
	var opts = metav1.ListOptions{
		LabelSelector: selector,
	}
	return client.
		CoreV1().
		Pods(cmd.conf.Namespace).
		List(ctx, opts)
}
