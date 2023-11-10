package cmd

import (
	"context"

	"github.com/enrichman/kubectl-epinio/internal/cli"
	"github.com/enrichman/kubectl-epinio/pkg/epinio"
	"github.com/spf13/cobra"
)

var emptyCompletions = []string{""}

type ValidArgsFunc func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)

// NoFileCompletions can be used to disable file completion for commands that should not trigger file completions.
func NoFileCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return emptyCompletions, cobra.ShellCompDirectiveNoFileComp
}

func NewUserValidator(epinioCLI *cli.EpinioCLI) ValidArgsFunc {
	return newGetterValidator(epinioCLI.KubeClient.ListUsers)
}

func NewRoleValidator(epinioCLI *cli.EpinioCLI) ValidArgsFunc {
	return newGetterValidator(epinioCLI.KubeClient.ListRoles)
}

func newGetterValidator[T idGetter](getter resourceGetter[T]) ValidArgsFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		resources, err := getResources(cmd.Context(), getter)
		if err != nil {
			return emptyCompletions, cobra.ShellCompDirectiveNoFileComp
		}

		return filter(args, resources)
	}
}

func NewNamespacesFlagValidator(epinioCLI *cli.EpinioCLI) ValidArgsFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		namespaces, err := epinioCLI.KubeClient.ListNamespaces(cmd.Context())
		if err != nil {
			return emptyCompletions, cobra.ShellCompDirectiveNoFileComp
		}

		alreadySelected, err := cmd.Flags().GetStringSlice("namespaces")
		if err != nil {
			return emptyCompletions, cobra.ShellCompDirectiveNoFileComp
		}

		return filter(alreadySelected, namespaces)
	}
}

func NewRolesFlagValidator(epinioCLI *cli.EpinioCLI) ValidArgsFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		roles, err := getResources(cmd.Context(), epinioCLI.KubeClient.ListRoles)
		if err != nil {
			return emptyCompletions, cobra.ShellCompDirectiveNoFileComp
		}

		alreadySelected, err := cmd.Flags().GetStringSlice("roles")
		if err != nil {
			return emptyCompletions, cobra.ShellCompDirectiveNoFileComp
		}

		return filter(alreadySelected, roles)
	}
}

func NewActionsFlagsValidator() ValidArgsFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		alreadySelected, err := cmd.Flags().GetStringSlice("actions")
		if err != nil {
			return emptyCompletions, cobra.ShellCompDirectiveNoFileComp
		}

		return filter(alreadySelected, epinio.Actions)
	}
}

type resourceGetter[T idGetter] func(context.Context) ([]T, error)

func getResources[T idGetter](ctx context.Context, getter resourceGetter[T]) ([]string, error) {
	resources, err := getter(ctx)
	if err != nil {
		return nil, err
	}

	return idGetterToStringArray(resources), nil
}

type idGetter interface {
	GetID() string
}

func filter(matched []string, resourceIDs []string) ([]string, cobra.ShellCompDirective) {
	// map to check for already selected resources
	alreadyMatched := map[string]struct{}{}
	for _, resource := range matched {
		alreadyMatched[resource] = struct{}{}
	}

	filtered := []string{}
	for _, id := range resourceIDs {
		if _, matched := alreadyMatched[id]; !matched {
			filtered = append(filtered, id)
		}
	}

	if len(filtered) > 0 {
		return filtered, cobra.ShellCompDirectiveNoFileComp
	}

	return emptyCompletions, cobra.ShellCompDirectiveNoFileComp
}

func idGetterToStringArray[T idGetter](resources []T) []string {
	ids := []string{}
	for _, r := range resources {
		ids = append(ids, r.GetID())
	}
	return ids
}
