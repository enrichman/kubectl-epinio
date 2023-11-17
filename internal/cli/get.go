package cli

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	_ "embed"

	"github.com/enrichman/kubectl-epinio/pkg/epinio"
)

func (e *EpinioCLI) GetRoles(ctx context.Context, names []string) error {
	roles, err := e.KubeClient.ListRoles(ctx)
	if err != nil {
		return err
	}

	roles = filterRoles(names, roles)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)

	fmt.Fprintln(w, "ID\tNAME\tDEFAULT\tACTIONS\tAGE")
	for _, r := range roles {
		age := "-"
		if r.CreationTimestamp.Unix() > 0 {
			age = time.Since(r.CreationTimestamp).Truncate(time.Second).String()
		}
		fmt.Fprintf(w, "%s\t%s\t%t\t%d\t%s\n", r.ID, r.Name, r.Default, len(r.Actions), age)
	}

	return w.Flush()
}

func filterRoles(usernames []string, roles []epinio.Role) []epinio.Role {
	if len(usernames) == 0 {
		return roles
	}

	usernamesMap := map[string]struct{}{}
	for _, u := range usernames {
		usernamesMap[u] = struct{}{}
	}

	filtered := []epinio.Role{}
	for _, role := range roles {
		if _, found := usernamesMap[role.ID]; found {
			filtered = append(filtered, role)
		}
	}
	return filtered
}
