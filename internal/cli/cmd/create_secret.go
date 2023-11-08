package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"syscall"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"

	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

func NewCreateUserCmd2(kubeClient kubernetes.Interface, opts *EpinioOptions) *cobra.Command {

	cmd := &cobra.Command{
		Use:          "secret",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {

			c.InOrStdin()

			username, err := ask("Username: ", opts.In)
			if err != nil {
				return err
			}

			hash, err := askPassword()
			if err != nil {
				return err
			}

			role, err := askRole()
			if err != nil {
				return err
			}

			workspaces, err := ask("Workspaces (space or comma separated): ", opts.In)
			if err != nil {
				return err
			}
			workspacesList := splitWorkspaces(workspaces)

			printUser(username, hash, role, workspacesList)
			userSecret := newUserSecret(username, hash, role, workspacesList)

			secretInterface := kubeClient.CoreV1().Secrets("epinio")
			sec, err := secretInterface.Create(context.Background(), userSecret, metav1.CreateOptions{})
			if err != nil {
				if !k8serrors.IsAlreadyExists(err) {
					return err
				}

				yesNo, err := ask("User already exists. Do you want to update it? [y/n] ", opts.In)
				if err != nil {
					return err
				}
				fmt.Print()

				yesNo = strings.ToLower(yesNo)
				if yesNo == "y" || yesNo == "yes" {
					sec, err := secretInterface.Update(context.Background(), userSecret, metav1.UpdateOptions{})
					if err != nil {
						return err
					}
					fmt.Printf("User '%s' updated\n", sec.Name)
				}
			} else {
				fmt.Printf("User '%s' created\n", sec.Name)
			}

			return nil
		},
	}

	return cmd
}

func ask(q string, in io.Reader) (string, error) {
	fmt.Print(q)

	reader := bufio.NewReader(in)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

func askPassword() ([]byte, error) {
	fmt.Print("Password: ")

	bytePassword1, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}

	fmt.Println("")
	fmt.Print("Retype password: ")

	bytePassword2, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}
	fmt.Println("")

	if string(bytePassword1) != string(bytePassword2) {
		return nil, errors.New("password doesn't match")
	}

	encrypted, err := bcrypt.GenerateFromPassword(bytePassword1, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return encrypted, nil
}

func askRole() (string, error) {
	prompt := promptui.Select{
		Label:        "Select a role",
		Items:        []string{"admin", "user"},
		HideHelp:     true,
		HideSelected: true,
	}

	_, role, err := prompt.Run()
	if err != nil {
		return "", err
	}

	fmt.Println("Role:", role)
	return role, nil
}

func splitWorkspaces(workspaces string) []string {
	workspacesMap := map[string]struct{}{}

	for _, w := range strings.Fields(workspaces) {
		for _, w2 := range strings.Split(w, ",") {
			w2 = strings.TrimSpace(w2)
			if w2 != "" {
				workspacesMap[w2] = struct{}{}
			}
		}
	}

	workspacesList := []string{}
	for k := range workspacesMap {
		workspacesList = append(workspacesList, k)
	}

	sort.Strings(workspacesList)
	return workspacesList
}

func printUser(username string, hash []byte, role string, workspaces []string) {
	fmt.Println("\nCreating user")

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendRow([]any{"Username", "username"})
	t.AppendRow([]any{"Password hash", string(hash)})
	t.AppendRow([]any{"Role", role})
	t.AppendRow([]any{"Workspaces", strings.Join(workspaces, "\n")})
	t.AppendSeparator()
	t.Render()

	fmt.Println()
}

func newUserSecret(username string, hash []byte, role string, workspaces []string) *v1.Secret {
	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      username,
			Namespace: "epinio",
			Labels: map[string]string{
				"foo":  "bar",
				"role": role,
			},
		},
		StringData: map[string]string{
			"username":   username,
			"password":   string(hash),
			"workspaces": strings.Join(workspaces, "\n"),
		},
		Type: "Opaque",
	}
}
