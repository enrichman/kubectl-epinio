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

		// map to check for already selected resources
		alreadyMatched := map[string]struct{}{}
		for _, resource := range args {
			alreadyMatched[resource] = struct{}{}
		}

		names := []string{}
		for _, u := range users {
			if _, matched := alreadyMatched[u.Username]; !matched {
				names = append(names, u.Username)
			}
		}

		return names, cobra.ShellCompDirectiveNoFileComp
	}
}
