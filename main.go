package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

type EpinioOptions struct {
	configFlags *genericclioptions.ConfigFlags
	genericiooptions.IOStreams
}

func main() {
	flags := pflag.NewFlagSet("kubectl-epinio", pflag.ExitOnError)
	pflag.CommandLine = flags

	streams := genericiooptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
	epinioCmd := NewCmdEpinio(streams)
	if err := epinioCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func NewCmdEpinio(streams genericiooptions.IOStreams) *cobra.Command {
	options := &EpinioOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		IOStreams:   streams,
	}

	cmd := &cobra.Command{
		Use:          "epinio something else [flags]",
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

	cmdCreate := &cobra.Command{
		Use: "create",
	}

	cmdCreate.AddCommand(
		NewCmdCreateUser(streams),
	)

	cmd.AddCommand(
		cmdCreate,
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
