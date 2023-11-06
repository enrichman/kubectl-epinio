package main

import (
	"os"

	"github.com/enrichman/kubectl-epinio/internal/cli/cmd"
	"github.com/spf13/pflag"

	"k8s.io/cli-runtime/pkg/genericiooptions"
)

func main() {
	flags := pflag.NewFlagSet("kubectl-epinio", pflag.ExitOnError)
	pflag.CommandLine = flags

	streams := genericiooptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}

	epinioCmd, err := cmd.NewRootCmd(streams)
	if err != nil {
		os.Exit(1)
	}

	if err := epinioCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
