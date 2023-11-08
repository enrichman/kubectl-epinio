package cli

import (
	"context"
	"errors"
	"fmt"
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

			namespaces, err = promptMultiSelect(epinioNamespaces)
			if err != nil {
				return err
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

	return e.KubeClient.CreateUser(ctx, user)
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
