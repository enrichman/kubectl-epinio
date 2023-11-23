package cli

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"
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

	roles = filterAndSortRoles(names, roles)

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

func filterAndSortRoles(ids []string, roles []epinio.Role) []epinio.Role {
	// if usernames are not specified sort and return all
	if len(ids) == 0 {
		slices.SortFunc(roles, func(i, j epinio.Role) int {
			return strings.Compare(i.ID, j.ID)
		})
		return roles
	}

	return filterRoles(ids, roles)
}

func filterRoles(ids []string, roles []epinio.Role) []epinio.Role {
	rolesMap := map[string]epinio.Role{}
	for _, role := range roles {
		rolesMap[role.ID] = role
	}

	filtered := []epinio.Role{}
	for _, id := range ids {
		if _, found := rolesMap[id]; found {
			filtered = append(filtered, rolesMap[id])
		}
	}
	return filtered
}
