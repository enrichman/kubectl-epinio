package cli

import (
	"context"
	"fmt"

	_ "embed"

	"github.com/enrichman/kubectl-epinio/pkg/epinio"
)

func (e *EpinioCLI) GetRoles(ctx context.Context, names []string) error {
	roles, err := e.KubeClient.ListRoles(ctx)
	if err != nil {
		return err
	}

	roles = filterRoles(names, roles)

	fmt.Println("ID")
	for _, r := range roles {
		fmt.Println(r.ID)
	}

	return nil
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
