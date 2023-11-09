package cmd

import (
	"github.com/enrichman/kubectl-epinio/internal/cli"
	"github.com/spf13/cobra"
)

type ValidArgsFunc func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)

func NewUserValidator(epinioCLI *cli.EpinioCLI) ValidArgsFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		users, err := epinioCLI.KubeClient.ListUsers(cmd.Context())
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		names := matchedFilter(args, users)
		return names, cobra.ShellCompDirectiveNoFileComp
	}
}

func NewRoleValidator(epinioCLI *cli.EpinioCLI) ValidArgsFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		roles, err := epinioCLI.KubeClient.ListRoles(cmd.Context())
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		names := matchedFilter(args, roles)
		return names, cobra.ShellCompDirectiveNoFileComp
	}
}

type idGetter interface {
	GetID() string
}

func matchedFilter[T idGetter](matched []string, resources []T) []string {
	// map to check for already selected resources
	alreadyMatched := map[string]struct{}{}
	for _, resource := range matched {
		alreadyMatched[resource] = struct{}{}
	}

	names := []string{}
	for _, r := range resources {
		if _, matched := alreadyMatched[r.GetID()]; !matched {
			names = append(names, r.GetID())
		}
	}

	return names
}
