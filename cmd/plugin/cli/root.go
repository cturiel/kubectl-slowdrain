package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/cturiel/kubectl-slowdrain/pkg/logger"
	"github.com/cturiel/kubectl-slowdrain/pkg/plugin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	KubernetesConfigFlags *genericclioptions.ConfigFlags
	delaySeconds          int
	infraPrefixes         []string
	autoConfirm           bool
	logLevel              string
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubectl-slowdrain <node_name>",
		Short: "Drains a node, deleting app pods one by one with delay.",
		Long:  "This plugin drains a Kubernetes node by deleting application pods one by one with a configurable delay, avoiding downtime for applications.",
		Args:  cobra.ExactArgs(1),
		RunE:  runSlowDrain,
	}

	cmd.Flags().IntVarP(&delaySeconds, "delay", "d", 20, "Seconds to wait between pod deletions")
	cmd.Flags().StringSliceVar(&infraPrefixes, "infra-prefixes", []string{"kube-", "openshift-", "infra-"}, "List of namespace prefixes for infra pods")
	cmd.Flags().BoolVarP(&autoConfirm, "assumeyes", "y", false, "Skip confirmation prompt")
	cmd.Flags().StringVar(&logLevel, "log-level", "info", "Logging level (debug, info, warn, error)")

	cmd.AddCommand(NewVersionCmd())

	KubernetesConfigFlags = genericclioptions.NewConfigFlags(true)
	KubernetesConfigFlags.AddFlags(cmd.Flags())

	cobra.OnInitialize(initConfig)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	return cmd
}

func InitAndExecute() {
	if err := RootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	viper.AutomaticEnv()
}

// Main function to run the drain of the node
func runSlowDrain(cmd *cobra.Command, args []string) error {
	log := logger.NewLogger(logLevel)
	nodeName := args[0]

	log.Debug("DEBUG log level enabled")
	log.Info("Initialize the node drain: %s", nodeName)

	config, err := KubernetesConfigFlags.ToRESTConfig()
	if err != nil {
		return fmt.Errorf("error getting Kubernetes configuration: %v", err)
	}

	clientset, err := plugin.NewKubeClient(config)
	if err != nil {
		return fmt.Errorf("error creating Kubernetes client: %v", err)
	}

	ctx := context.Background()
	return plugin.DrainNode(ctx, clientset, nodeName, delaySeconds, autoConfirm, logLevel, infraPrefixes)
}
