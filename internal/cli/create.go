package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/AlecAivazis/survey/v2"
	"github.com/enrichman/kubectl-epinio/pkg/epinio"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

func (e *EpinioCLI) CreateUser(ctx context.Context, username, password string, namespaces, roles []string, interactive bool) error {
	var err error

	// start interactive mode
	if interactive {
		if password == "" {
			password, err = promptPassword()
			if err != nil {
				return err
			}
		}

		if len(namespaces) == 0 {
			epinioNamespaces, err := e.KubeClient.ListNamespaces(ctx)
			if err != nil {
				return err
			}

			msg := "Namespaces assigned to the user:"
			namespaces, err = promptMultiSelect(msg, epinioNamespaces)
			if err != nil {
				return err
			}
		}

		if len(roles) == 0 {
			epinioRoles, err := e.KubeClient.ListRoles(ctx)
			if err != nil {
				return err
			}

			roleIDs := []string{}
			for _, r := range epinioRoles {
				roleIDs = append(roleIDs, r.GetID())
			}

			msg := "Global Roles assigned to the user:"
			roles, err = promptMultiSelect(msg, roleIDs)
			if err != nil {
				return err
			}

			// you need to have some namespaces assigned for Namescoped Roles
			if len(namespaces) > 0 {
				confirm, err := promptConfirmation("Do you want to assign Namescoped Roles? [y/n] ")
				if err != nil {
					return err
				}
				if confirm {
					var namescopedRoles []string
					msg := "Namescoped Roles assigned to the user:"

					var selectedRoles []string
					selectedRoles, err = promptMultiSelect(msg, roleIDs)
					if err != nil {
						return err
					}

					for _, roleToNamescope := range selectedRoles {
						msg := fmt.Sprintf("Namespaces for '%s' role:", roleToNamescope)
						namespacesForRole, err := promptMultiSelect(msg, namespaces)
						if err != nil {
							return err
						}

						for _, ns := range namespacesForRole {
							namescopedRole := fmt.Sprintf("%s:%s", roleToNamescope, ns)
							namescopedRoles = append(namescopedRoles, namescopedRole)
						}
					}

					roles = append(roles, namescopedRoles...)
				}
			}
		}
	}

	password, err = hashPasswordIfNeeded(password)
	if err != nil {
		return err
	}

	user := epinio.User{
		Username:   username,
		Password:   password,
		Namespaces: namespaces,
		Roles:      roles,
	}

	fmt.Println()
	fmt.Printf(format, "Username:", user.Username)
	fmt.Printf(format, "Password:", user.Password)
	printArray("Roles:", user.Roles)
	printArray("Namespaces:", user.Namespaces)
	fmt.Println()

	if interactive {
		confirm, err := promptConfirmation("Create? [y/n] ")
		if err != nil {
			return err
		}

		if !confirm {
			fmt.Println("aborted!")
			return nil
		}
	}

	err = e.KubeClient.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	fmt.Println("User created!")
	return nil
}

func (e *EpinioCLI) CreateRole(ctx context.Context, id, name string, isDefault bool, actions []string, interactive bool) error {
	var err error

	// start interactive mode
	if interactive {
		fmt.Println("ID:", id)

		if name == "" {
			name, err = promptName()
			if err != nil {
				return err
			}
		}

		if !isDefault {
			isDefault, err = promptDefault()
			if err != nil {
				return err
			}
		}

		if len(actions) == 0 {
			msg := "Actions assigned to the role:"
			actions, err = promptMultiSelect(msg, epinio.Actions)
			if err != nil {
				return err
			}
		}

		//TODO: print summary (also in non-interactive mode?)

		confirm, err := promptConfirmation("Create? [y/n] ")
		if err != nil {
			return err
		}

		if !confirm {
			fmt.Println("aborted!")
			return nil
		}
	}

	role := epinio.Role{
		ID:      id,
		Name:    name,
		Default: isDefault,
		Actions: actions,
	}

	err = e.KubeClient.CreateRole(ctx, role)
	if err != nil {
		return err
	}

	fmt.Println("Role created!")
	return nil
}

func promptPassword() (string, error) {
	fmt.Print("Password: ")

	bytePassword1, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	fmt.Print("\nRetype password: ")

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

func promptName() (string, error) {
	fmt.Print("Name: ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(input), nil
}

func promptDefault() (bool, error) {
	fmt.Print("Is Default? [true/false] ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	boolStr := strings.TrimSpace(input)

	isDefault, _ := strconv.ParseBool(boolStr)
	return isDefault, nil
}

func promptConfirmation(msg string) (bool, error) {
	fmt.Print(msg)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	yn := strings.TrimSpace(input)
	yn = strings.ToLower(yn)

	return (yn == "y" || yn == "yes"), nil
}

func promptMultiSelect(msg string, options []string) ([]string, error) {
	prompt := &survey.MultiSelect{
		Message:  msg,
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

func hashPasswordIfNeeded(password string) (string, error) {
	if password == "" {
		return password, nil
	}

	if len(password) == 60 && strings.HasPrefix(password, "$2a$") {
		return password, nil
	}

	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(encrypted), nil
}
