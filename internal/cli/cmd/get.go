package cmd

import (
	"github.com/enrichman/kubectl-epinio/internal/cli"
	"github.com/spf13/cobra"
)

func NewGetCmd(cli *cli.EpinioCLI) *cobra.Command {
	getCmd := &cobra.Command{
		Use: "get",
		RunE: func(c *cobra.Command, args []string) error {
			return nil
		},
	}

	getCmd.AddCommand(
		NewGetUserCmd(cli),
	)

	return getCmd
}

func NewGetUserCmd(cli *cli.EpinioCLI) *cobra.Command {
	return &cobra.Command{
		Use:               "user",
		Aliases:           []string{"users"},
		ValidArgsFunction: NewUserValidator(cli),
		RunE: func(c *cobra.Command, args []string) error {
			usernames := args
			return cli.GetUsers(c.Context(), usernames)
		},
	}
}
