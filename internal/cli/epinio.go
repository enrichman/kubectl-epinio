package cli

import (
	"context"
	"fmt"

	"github.com/enrichman/kubectl-epinio/pkg/epinio"
	"k8s.io/client-go/kubernetes"

	_ "embed"
)

// EpinioCLI handles the CLI commands, calling the Kubernetes API and handling the display
type EpinioCLI struct {
	KubeClient *epinio.KubeClient
	//TODO: add output and logs
}

func NewEpinioCLI(kubeClient kubernetes.Interface) (e *EpinioCLI) {
	return &EpinioCLI{
		KubeClient: epinio.NewKubeClient(kubeClient),
	}
}

func (e *EpinioCLI) GetUsers(ctx context.Context, usernames []string) error {
	users, err := e.KubeClient.ListUsers(ctx)
	if err != nil {
		return err
	}

	users = filterUsers(usernames, users)

	fmt.Println("USERNAME")
	for _, u := range users {
		fmt.Println(u.Username)
	}

	return nil
}

func (e *EpinioCLI) DescribeUsers(ctx context.Context, usernames []string) error {
	users, err := e.KubeClient.ListUsers(ctx)
	if err != nil {
		return err
	}

	users = filterUsers(usernames, users)

	format := "%-10s %s\n"

	for i, u := range users {
		if i != 0 {
			fmt.Println()
		}
		fmt.Printf(format, "Username:", u.Username)
		fmt.Printf(format, "Password:", u.Password)
	}

	return nil
}

func filterUsers(usernames []string, users []epinio.User) []epinio.User {
	if len(usernames) == 0 {
		return users
	}

	usernamesMap := map[string]struct{}{}
	for _, u := range usernames {
		usernamesMap[u] = struct{}{}
	}

	filtered := []epinio.User{}
	for _, user := range users {
		if _, found := usernamesMap[user.Username]; found {
			filtered = append(filtered, user)
		}
	}
	return filtered
}
