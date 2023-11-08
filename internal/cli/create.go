package cli

import (
	"context"

	_ "embed"
)

func (e *EpinioCLI) CreateUser(ctx context.Context, username, password string, namespaces, roles []string) error {
	return nil
}
