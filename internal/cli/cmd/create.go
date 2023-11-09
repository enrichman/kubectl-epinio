package cmd

import (
	"github.com/enrichman/kubectl-epinio/internal/cli"
	"github.com/spf13/cobra"
)

func NewCreateCmd(epinioCLI *cli.EpinioCLI) *cobra.Command {
	createCmd := &cobra.Command{
		Use: "create",
	}

	createCmd.AddCommand(
		NewCreateUserCmd(epinioCLI),
	)

	return createCmd
}

type CreateUserConfig struct {
	Interactive bool
	Password    string
	Namespaces  []string
	Roles       []string
}

func NewCreateUserCmd(epinioCLI *cli.EpinioCLI) *cobra.Command {
	cfg := &CreateUserConfig{}

	createUserCmd := &cobra.Command{
		Use:     "user <username>",
		Short:   "username/email of the user used during the login",
		Example: `kubectl epinio create user "foo@bar.io"`,
		Args:    cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			ctx := c.Context()
			username := args[0]

			namespaces := unique(cfg.Namespaces)

			return epinioCLI.CreateUser(ctx, username, cfg.Password, namespaces, cfg.Roles, cfg.Interactive)
		},
	}

	createUserCmd.Flags().BoolVarP(&cfg.Interactive, "interactive", "i", false, "interactive mode")
	createUserCmd.Flags().StringVar(&cfg.Password, "password", "", "plain password of the user used during the login")
	createUserCmd.Flags().StringSliceVar(&cfg.Namespaces, "namespaces", nil, "namespaces")
	createUserCmd.Flags().StringSliceVar(&cfg.Roles, "roles", nil, "roles")

	err := createUserCmd.RegisterFlagCompletionFunc("namespaces", NewNamespaceValidator(epinioCLI))
	checkErr(err, "cannot create 'create user' command")

	return createUserCmd
}
