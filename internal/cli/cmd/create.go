package cmd

import (
	"github.com/enrichman/kubectl-epinio/internal/cli"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

func NewCreateCmd(epinioCLI *cli.EpinioCLI) *cobra.Command {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create an Epinio resource [user/role]",
	}

	createCmd.AddCommand(
		NewCreateUserCmd(epinioCLI),
		NewCreateRoleCmd(epinioCLI),
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
		Use:               "user <username>",
		Short:             "Create a user",
		Example:           `kubectl epinio create user "foo@bar.io"`,
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: NoFileCompletions,
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

	err := createUserCmd.RegisterFlagCompletionFunc("namespaces", NewNamespacesFlagValidator(epinioCLI))
	checkErr(err, "cannot create 'create user' command")

	err = createUserCmd.RegisterFlagCompletionFunc("roles", NewRolesFlagValidator(epinioCLI))
	checkErr(err, "cannot create 'create user' command")

	return createUserCmd
}

type CreateRoleConfig struct {
	Interactive bool
	ID          string
	Name        string
	Default     bool
	Actions     []string
}

func NewCreateRoleCmd(epinioCLI *cli.EpinioCLI) *cobra.Command {
	cfg := &CreateRoleConfig{}

	createRoleCmd := &cobra.Command{
		Use:               "role <role_id>",
		Short:             "Create a role",
		Example:           `kubectl epinio create role "read_role"`,
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: NoFileCompletions,
		RunE: func(c *cobra.Command, args []string) error {
			ctx := c.Context()
			roleID := args[0]
			actions := unique(cfg.Actions)

			return epinioCLI.CreateRole(ctx, roleID, cfg.Name, cfg.Default, actions, cfg.Interactive)
		},
	}

	createRoleCmd.Flags().BoolVarP(&cfg.Interactive, "interactive", "i", false, "interactive mode")
	createRoleCmd.Flags().StringVar(&cfg.Name, "name", "", "friendly name of the role")
	createRoleCmd.Flags().BoolVar(&cfg.Default, "default", false, "set the role as default")
	createRoleCmd.Flags().StringSliceVar(&cfg.Actions, "actions", nil, "actions allowed for the role")

	err := createRoleCmd.RegisterFlagCompletionFunc("actions", NewActionsFlagsValidator())
	checkErr(err, "cannot create 'create role' command")

	return createRoleCmd
}

func unique(arr []string) []string {
	unique := map[string]struct{}{}

	for _, v := range arr {
		unique[v] = struct{}{}
	}

	return maps.Keys(unique)
}
