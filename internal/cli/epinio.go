package cli

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

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

	users = filterAndSortUsers(usernames, users)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)

	fmt.Fprintln(w, "USERNAME\tADMIN\tROLES\tAGE")
	for _, u := range users {
		var roles int
		if u.Role != "" {
			roles = 1
		} else if len(u.Roles) > 0 {
			roles = len(u.Roles)
		}

		age := "-"
		if u.CreationTimestamp.Unix() > 0 {
			age = time.Since(u.CreationTimestamp).Truncate(time.Second).String()
		}

		fmt.Fprintf(w, "%s\t%t\t%d\t%s\n", u.Username, u.IsAdmin(), roles, age)
	}

	return w.Flush()
}

func (e *EpinioCLI) DescribeUsers(ctx context.Context, usernames []string) error {
	users, err := e.KubeClient.ListUsers(ctx)
	if err != nil {
		return err
	}

	users = filterAndSortUsers(usernames, users)

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

	roles = filterAndSortRoles(ids, roles)

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

func filterAndSortUsers(usernames []string, users []epinio.User) []epinio.User {
	// if usernames are not specified sort and return all
	if len(usernames) == 0 {
		slices.SortFunc(users, func(i, j epinio.User) int {
			return strings.Compare(i.Username, j.Username)
		})
		return users
	}

	return filterUsers(usernames, users)
}

func filterUsers(usernames []string, users []epinio.User) []epinio.User {
	usersMap := map[string]epinio.User{}
	for _, u := range users {
		usersMap[u.Username] = u
	}

	filtered := []epinio.User{}
	for _, username := range usernames {
		if _, found := usersMap[username]; found {
			filtered = append(filtered, usersMap[username])
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
