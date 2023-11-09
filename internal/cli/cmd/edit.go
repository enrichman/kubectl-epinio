package cmd

import (
	"github.com/enrichman/kubectl-epinio/internal/cli"
	"github.com/spf13/cobra"
)

func NewEditCmd(cli *cli.EpinioCLI) *cobra.Command {
	editCmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit a resource from the default editor",
	}

	editCmd.AddCommand(
		NewEditUserCmd(cli),
	)

	return editCmd
}

func NewEditUserCmd(cli *cli.EpinioCLI) *cobra.Command {
	return &cobra.Command{
		Use:               "user",
		Short:             "Edit a user resource from the default editor",
		Aliases:           []string{"users"},
		ValidArgsFunction: NewUserValidator(cli),
		Args:              cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return cli.EditUser(c.Context(), args[0])
		},
	}
}
