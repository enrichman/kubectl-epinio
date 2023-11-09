package cmd

import (
	"github.com/enrichman/kubectl-epinio/internal/cli"
	"github.com/spf13/cobra"
)

func NewGetCmd(cli *cli.EpinioCLI) *cobra.Command {
	getCmd := &cobra.Command{
		Use: "get",
	}

	getCmd.AddCommand(
		NewGetUserCmd(cli),
		NewGetRoleCmd(cli),
	)

	return getCmd
}

func NewGetUserCmd(cli *cli.EpinioCLI) *cobra.Command {
	return &cobra.Command{
		Use:               "user",
		Aliases:           []string{"users"},
		ValidArgsFunction: NewUserValidator(cli),
		RunE: func(c *cobra.Command, args []string) error {
			return cli.GetUsers(c.Context(), args)
		},
	}
}

func NewGetRoleCmd(cli *cli.EpinioCLI) *cobra.Command {
	return &cobra.Command{
		Use:               "role",
		Aliases:           []string{"roles"},
		ValidArgsFunction: NewRoleValidator(cli),
		RunE: func(c *cobra.Command, args []string) error {
			return cli.GetRoles(c.Context(), args)
		},
	}
}
