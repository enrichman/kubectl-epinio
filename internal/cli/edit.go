package cli

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/template"

	"github.com/enrichman/kubectl-epinio/pkg/epinio"
	"gopkg.in/yaml.v3"

	_ "embed"
)

//go:embed template/edit_user.yaml
var editUserTemplate string

func (e *EpinioCLI) EditUser(ctx context.Context, username string) error {
	user, err := e.KubeClient.GetUser(ctx, username)
	if err != nil {
		return err
	}

	tempFile, err := createEditUserTempFile(user)
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())

	err = openEditor(tempFile)
	if err != nil {
		return err
	}

	tempFile, err = os.Open(tempFile.Name())
	if err != nil {
		return err
	}

	updatedUser := &epinio.User{}
	err = yaml.NewDecoder(tempFile).Decode(updatedUser)
	if err != nil {
		return err
	}

	fmt.Println("updated", updatedUser)

	return nil
}

func createEditUserTempFile(user epinio.User) (*os.File, error) {
	tempFilename := fmt.Sprintf("%s/edit-user-%s.yaml", os.TempDir(), user.Username)

	tempFile, err := os.Create(tempFilename)
	if err != nil {
		return nil, err
	}

	templ, err := template.New("editUserTemplate").Parse(editUserTemplate)
	if err != nil {
		return nil, err
	}

	err = templ.Execute(tempFile, user)
	if err != nil {
		return nil, err
	}
	tempFile.Close()

	return tempFile, nil
}

func openEditor(tempFile *os.File) error {
	defaultEditor := "vim"
	if editor, found := os.LookupEnv("EDITOR"); found {
		defaultEditor = editor
	}

	cmd := exec.Command(defaultEditor, tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("Error while editing. Error: %v\n", err)
		return err
	}

	log.Printf("Successfully edited.")
	return nil
}
