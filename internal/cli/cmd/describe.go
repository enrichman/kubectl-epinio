package cmd

import (
	"github.com/enrichman/kubectl-epinio/internal/cli"
	"github.com/spf13/cobra"
)

func NewDescribeCmd(cli *cli.EpinioCLI) *cobra.Command {
	describeCmd := &cobra.Command{
		Use: "describe",
		RunE: func(c *cobra.Command, args []string) error {
			return nil
		},
	}

	describeCmd.AddCommand(
		NewDescribeUserCmd(cli),
	)

	return describeCmd
}

func NewDescribeUserCmd(cli *cli.EpinioCLI) *cobra.Command {
	return &cobra.Command{
		Use:     "user",
		Aliases: []string{"users"},
		RunE: func(c *cobra.Command, args []string) error {
			return cli.Get(c.Context(), args[0])
		},
	}
}
