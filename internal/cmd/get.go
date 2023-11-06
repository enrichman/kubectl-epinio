package cmd

import (
	"github.com/spf13/cobra"
)

func NewGetCmd() *cobra.Command {
	getCmd := &cobra.Command{
		Use: "get",
		RunE: func(c *cobra.Command, args []string) error {
			return nil
		},
	}

	getCmd.AddCommand(
		NewGetUserCmd(),
	)

	return getCmd
}

func NewGetUserCmd() *cobra.Command {
	return &cobra.Command{
		Use: "user",
		RunE: func(c *cobra.Command, args []string) error {
			return nil
		},
	}
}
