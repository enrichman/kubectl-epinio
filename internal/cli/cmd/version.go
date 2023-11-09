package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
)

var Version = "v0.0.0-dev"

func NewVersionCmd(kubeClient kubernetes.Interface) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the plugin and server version",
		RunE: func(c *cobra.Command, args []string) error {
			sv, err := kubeClient.Discovery().ServerVersion()
			if err != nil {
				return err
			}

			fmt.Printf("kubectl-epinio Version: %s\n", Version)
			fmt.Printf("Server Version: %s\n", sv.String())

			return nil
		},
	}
}
