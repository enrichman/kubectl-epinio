package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/enrichman/kubectl-epinio/pkg/epinio"
	"github.com/liggitt/tabwriter"
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

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)
	fmt.Fprintln(w, "USERNAME\tADMIN\tROLES\tAGE")
	for _, u := range users {
		isAdmin := isUserAdmin(u)
		fmt.Fprintf(w, "%s\t%t\t%d\t%s\n", u.Username, isAdmin, len(u.Roles), time.Since(u.CreationTimestamp).String())
	}

	return w.Flush()
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
			fmt.Println()
		}
		fmt.Printf(format, "Username:", u.Username)
		fmt.Printf(format, "Password:", u.Password)
		printArray("Roles:", u.Roles)
		printArray("Namespaces:", u.Namespaces)
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
}

func isUserAdmin(user epinio.User) bool {
	if user.Role == "admin" {
		return true
	}
	for _, r := range user.Roles {
		if r == "admin" {
			return true
		}
	}
	return false
}
