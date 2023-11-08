package cmd

import (
	"errors"
	"fmt"
	"strings"
	"syscall"

	"github.com/enrichman/kubectl-epinio/internal/cli"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"

	"github.com/AlecAivazis/survey/v2"
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
			username := args[0]

			// start interactive mode
			if cfg.Interactive {
				if err := interactiveMode(cfg); err != nil {
					return err
				}
			}

			fmt.Println("pp", cfg.Password)
			password, err := hashPasswordIfNeeded(cfg.Password)
			if err != nil {
				return err
			}
			cfg.Password = password

			fmt.Println(cfg)

			return epinioCLI.CreateUser(c.Context(), username, cfg.Password, cfg.Namespaces, cfg.Roles)
		},
	}

	createUserCmd.Flags().BoolVarP(&cfg.Interactive, "interactive", "i", false, "interactive mode")
	createUserCmd.Flags().StringVar(&cfg.Password, "password", "", "plain password of the user used during the login")
	createUserCmd.Flags().StringSliceVar(&cfg.Namespaces, "namespace", nil, "namespaces")
	createUserCmd.Flags().StringSliceVar(&cfg.Roles, "role", nil, "roles")

	return createUserCmd
}

func interactiveMode(cfg *CreateUserConfig) error {
	if cfg.Password == "" {
		password, err := promptPassword()
		if err != nil {
			return err
		}

		cfg.Password = password
	}

	if len(cfg.Namespaces) == 0 {
		availableNamespaces := []string{"workspace", "workspace2"}
		namespaces, err := promptMultiSelect(availableNamespaces)
		if err != nil {
			return err
		}
		cfg.Namespaces = namespaces
	}

	return nil
}

func promptPassword() (string, error) {
	fmt.Print("Password: ")

	bytePassword1, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	fmt.Println("")
	fmt.Print("Retype password: ")

	bytePassword2, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Println("")

	if string(bytePassword1) != string(bytePassword2) {
		return "", errors.New("password doesn't match")
	}

	return string(bytePassword1), nil
}

func promptMultiSelect(options []string) ([]string, error) {
	prompt := &survey.MultiSelect{
		Message:  "Namespaces assigned to the user:",
		Options:  options,
		PageSize: 10,
	}

	selected := []string{}
	err := survey.AskOne(
		prompt,
		&selected,
		survey.WithRemoveSelectAll(),
		survey.WithRemoveSelectNone(),
	)
	if err != nil {
		return nil, err
	}
	return selected, nil
}

// $2a$10$uD1k7WFEnllu26uDFXNsGuJmGuyiF.nfNSgo.s41fVb8ce9eKSj.6

func hashPasswordIfNeeded(password string) (string, error) {
	if len(password) == 60 && strings.HasPrefix(password, "$2a$") {
		return password, nil
	}

	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(encrypted), nil
}
