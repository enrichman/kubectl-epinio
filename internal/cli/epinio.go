package cli

import (
	"context"
	"fmt"

	"github.com/enrichman/kubectl-epinio/pkg/epinio"
	"k8s.io/client-go/kubernetes"
)

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

func (e *EpinioCLI) Get(ctx context.Context, username string) error {
	users, err := e.KubeClient.ListUsers(ctx)
	if err != nil {
		return err
	}

	fmt.Println("USERNAME")
	for _, u := range users {
		fmt.Println(u.Username)
	}

	return nil
}

func (e *EpinioCLI) Describe(ctx context.Context, username string) error {
	users, err := e.KubeClient.ListUsers(ctx)
	if err != nil {
		return err
	}

	format := "%-10s %s\n"

	for _, u := range users {
		fmt.Printf(format, "Username:", u.Username)
		fmt.Printf(format, "Password:", u.Password)
		fmt.Println()
	}

	return nil
}
