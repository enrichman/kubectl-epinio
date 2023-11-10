package cmd

import (
	"github.com/enrichman/kubectl-epinio/internal/cli"
	"github.com/spf13/cobra"
)

func NewDeleteCmd(epinioCLI *cli.EpinioCLI) *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an Epinio resource [user/role]",
	}

	deleteCmd.AddCommand(
		NewDeleteUserCmd(epinioCLI),
		NewDeleteRoleCmd(epinioCLI),
	)

	return deleteCmd
}

func NewDeleteUserCmd(epinioCLI *cli.EpinioCLI) *cobra.Command {
	var noConfirm bool

	deleteUserCmd := &cobra.Command{
		Use:               "user <username>",
		Short:             "Delete a user",
		Example:           `kubectl epinio delete user "foo@bar.io"`,
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: NoFileCompletions,
		RunE: func(c *cobra.Command, args []string) error {
			ctx := c.Context()
			username := args[0]

			return epinioCLI.DeleteUser(ctx, username, noConfirm)
		},
	}

	deleteUserCmd.Flags().BoolVarP(&noConfirm, "no-confirm", "y", false, "delete without confirmation")

	return deleteUserCmd
}

func NewDeleteRoleCmd(epinioCLI *cli.EpinioCLI) *cobra.Command {
	var noConfirm bool

	deleteRoleCmd := &cobra.Command{
		Use:               "role <role_id>",
		Short:             "Delete a role",
		Example:           `kubectl epinio delete role "read_role"`,
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: NoFileCompletions,
		RunE: func(c *cobra.Command, args []string) error {
			ctx := c.Context()
			id := args[0]

			return epinioCLI.DeleteRole(ctx, id, noConfirm)
		},
	}

	deleteRoleCmd.Flags().BoolVarP(&noConfirm, "no-confirm", "y", false, "deletion without confirmation")

	return deleteRoleCmd
}
