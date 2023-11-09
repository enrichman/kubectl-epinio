package cmd

import (
	"strings"

	"github.com/enrichman/kubectl-epinio/internal/cli"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

type ValidArgsFunc func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)

// NoFileCompletions can be used to disable file completion for commands that should not trigger file completions.
func NoFileCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{""}, cobra.ShellCompDirectiveNoFileComp
}

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

func NewNamespaceValidator(epinioCLI *cli.EpinioCLI) ValidArgsFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		ns, err := epinioCLI.KubeClient.ListNamespaces(cmd.Context())
		if err != nil {
			return []string{""}, cobra.ShellCompDirectiveNoFileComp
		}

		alreadySelected, err := cmd.Flags().GetStringSlice("namespaces")
		if err != nil {
			return []string{""}, cobra.ShellCompDirectiveNoFileComp
		}

		filtered := filter(alreadySelected, ns)
		if len(filtered) > 0 {
			return filtered, cobra.ShellCompDirectiveNoFileComp
		}

		return []string{""}, cobra.ShellCompDirectiveNoFileComp
	}
}

func NewStaticFlagsCompletionFunc(allowedValues []string) ValidArgsFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		matches := []string{}

		for _, allowed := range allowedValues {
			if strings.HasPrefix(allowed, toComplete) {
				matches = append(matches, allowed)
			}
		}

		alreadySelected, err := cmd.Flags().GetStringSlice("actions")
		if err != nil {
			return []string{""}, cobra.ShellCompDirectiveNoFileComp
		}

		filtered := filter(alreadySelected, matches)
		if len(filtered) > 0 {
			return filtered, cobra.ShellCompDirectiveNoFileComp
		}

		return matches, cobra.ShellCompDirectiveNoFileComp
	}
}

type idGetter interface {
	GetID() string
}

func matchedFilter[T idGetter](matched []string, resources []T) []string {
	converted := idGetterToStringArray(resources)
	return filter(matched, converted)
}

func filter(matched []string, resourceIDs []string) []string {
	// map to check for already selected resources
	alreadyMatched := map[string]struct{}{}
	for _, resource := range matched {
		alreadyMatched[resource] = struct{}{}
	}

	names := []string{}
	for _, id := range resourceIDs {
		if _, matched := alreadyMatched[id]; !matched {
			names = append(names, id)
		}
	}

	return names
}

func idGetterToStringArray[T idGetter](resources []T) []string {
	ids := []string{}
	for _, r := range resources {
		ids = append(ids, r.GetID())
	}
	return ids
}

func unique[T comparable](arr []T) []T {
	unique := map[T]struct{}{}

	for _, v := range arr {
		unique[v] = struct{}{}
	}

	return maps.Keys(unique)
}
