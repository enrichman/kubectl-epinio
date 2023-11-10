package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/enrichman/kubectl-epinio/pkg/epinio"
	"k8s.io/client-go/kubernetes"

	_ "embed"
)

const format = "%-15s %s\n"

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

	format := "%-15s %s\n"

	for i, u := range users {
		if i != 0 {
			fmt.Println("---")
		}
		fmt.Printf(format, "Username:", u.Username)
		fmt.Printf(format, "Password:", u.Password)
		printArray("Roles:", u.Roles)
		printArray("Namespaces:", u.Namespaces)
	}

	return nil
}

func (e *EpinioCLI) DescribeRoles(ctx context.Context, ids []string) error {
	roles, err := e.KubeClient.ListRoles(ctx)
	if err != nil {
		return err
	}

	roles = filterRoles(ids, roles)

	format := "%-15s %s\n"

	for i, r := range roles {
		if i != 0 {
			fmt.Println("---")
		}
		fmt.Printf(format, "ID:", r.ID)
		fmt.Printf(format, "Name:", r.Name)
		fmt.Printf(format, "Default:", strconv.FormatBool(r.Default))
		printArray("Actions:", r.Actions)
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

func printArray(description string, arr []string) {
	if len(arr) == 0 {
		fmt.Printf(format, description, "")
	}

	for i, ns := range arr {
		var leftCol string
		if i == 0 {
			leftCol = description
		}
		fmt.Printf(format, leftCol, ns)
	}
	fmt.Println()
}
