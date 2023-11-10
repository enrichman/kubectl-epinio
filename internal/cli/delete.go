package cli

import (
	"context"
	"fmt"
)

// DeleteUser deletes a user by username, after asking for confirmation if noConfirm resolves to `false`.
func (e *EpinioCLI) DeleteUser(ctx context.Context, username string, noConfirm bool) error {
	if !noConfirm {
		confirm, err := promptConfirmation("Delete? [y/n] ")
		if err != nil {
			return err
		}

		if !confirm {
			fmt.Println("aborted!")
			return nil
		}
	}

	if err := e.KubeClient.DeleteUser(ctx, username); err != nil {
		return err
	}

	fmt.Println("User deleted!")

	return nil
}

// DeleteRole deletes a role by id, after asking for confirmation if noConfirm resolves to `false`.
func (e *EpinioCLI) DeleteRole(ctx context.Context, id string, noConfirm bool) error {
	if !noConfirm {
		confirm, err := promptConfirmation("Delete? [y/n] ")
		if err != nil {
			return err
		}

		if !confirm {
			fmt.Println("aborted!")
			return nil
		}
	}

	if err := e.KubeClient.DeleteRole(ctx, id); err != nil {
		return err
	}

	fmt.Println("Role deleted!")

	return nil
}
