package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type EpinioOptions struct {
	configFlags *genericclioptions.ConfigFlags
	genericiooptions.IOStreams
}

func NewEpinioCmd(streams genericiooptions.IOStreams) *cobra.Command {
	options := &EpinioOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		IOStreams:   streams,
	}

	cmd := &cobra.Command{
		Use:          "epinio",
		Short:        "blabla epinio",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			config, err := options.configFlags.ToRESTConfig()
			if err != nil {
				return err
			}

			discoveryClient, err := options.configFlags.ToDiscoveryClient()
			if err != nil {
				return err
			}

			_, err = discoveryClient.ServerVersion()
			if err != nil {
				return err
			}

			version, err := getServerVersion(config)
			if err != nil {
				return err
			}

			fmt.Println(version)

			return nil
		},
	}

	cmd.AddCommand(
		NewVersionCmd(),
		NewCreateUserCmd(streams),
		NewGetCmd(),
	)

	options.configFlags.AddFlags(cmd.Flags())

	return cmd
}

func getServerVersion(config *rest.Config) (string, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}

	sv, err := clientset.Discovery().ServerVersion()
	if err != nil {
		return "", err
	}

	return sv.String(), nil
}
