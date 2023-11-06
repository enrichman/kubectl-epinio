package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "v0.0.0-dev"

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use: "version",
		Run: func(c *cobra.Command, args []string) {
			fmt.Printf("kubectl-epinio Version: %s\n", Version)
		},
	}
}
