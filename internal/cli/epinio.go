package cli

import (
	"context"
	"fmt"

	"github.com/enrichman/kubectl-epinio/pkg/epinio"
	"k8s.io/client-go/kubernetes"
)

// EpinioCLI handles the CLI commands, calling the Kubernetes API and handling the display
type EpinioCLI struct {
	KubeClient kubernetes.Interface
}

func NewEpinioCLI(kubeClient kubernetes.Interface) (e *EpinioCLI) {
	return &EpinioCLI{
		KubeClient: kubeClient,
	}
}

func (e *EpinioCLI) Get(username string) error {
	users, err := epinio.ListUsers(context.Background(), e.KubeClient)
	if err != nil {
		return err
	}

	fmt.Println("USERNAME")
	for _, u := range users {
		fmt.Println(u.Username)
	}

	return nil
}

func (e *EpinioCLI) Describe(username string) error {
	users, err := epinio.ListUsers(context.Background(), e.KubeClient)
	if err != nil {
		return err
	}

	for _, u := range users {
		fmt.Printf("%-10s %s\n", "Username:", u.Username)
		fmt.Printf("%-10s %s\n", "Password:", u.Password)
		fmt.Println()
	}

	return nil
}
