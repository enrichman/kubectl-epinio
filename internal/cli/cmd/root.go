package cmd

import (
	"log"

	"github.com/enrichman/kubectl-epinio/internal/cli"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/client-go/kubernetes"
)

type EpinioOptions struct {
	configFlags *genericclioptions.ConfigFlags
	genericiooptions.IOStreams
}

func NewRootCmd(streams genericiooptions.IOStreams) (*cobra.Command, error) {
	options := &EpinioOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		IOStreams:   streams,
	}

	config, err := options.configFlags.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	epinioCLI := cli.NewEpinioCLI(kubeClient)

	rootCmd := &cobra.Command{
		Use:          "epinio",
		Short:        "Kubectl plugin to manage some Epinio resources",
		SilenceUsage: true,
	}

	rootCmd.AddCommand(
		NewVersionCmd(kubeClient),
		NewGetCmd(epinioCLI),
		NewDescribeCmd(epinioCLI),
		NewEditCmd(epinioCLI),
		NewCreateCmd(epinioCLI),
	)

	// Uncomment to add some kubectl flags
	// options.configFlags.AddFlags(rootCmd.Flags())

	return rootCmd, nil
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatal(msg, err)
	}
}
