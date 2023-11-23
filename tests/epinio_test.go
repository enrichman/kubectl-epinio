package tests

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
)

type KubectlEpinio struct {
	binPath string
}

func NewKubectlEpinio() (*KubectlEpinio, error) {
	cmd := exec.Command(
		"git",
		"rev-parse",
		"--show-toplevel",
	)

	root, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	repoRoot := strings.TrimSpace(string(root))

	return &KubectlEpinio{
		binPath: path.Join(repoRoot, "output", "kubectl-epinio"),
	}, nil
}

func (e *KubectlEpinio) Run(args ...string) (string, string, error) {
	cmd := exec.Command(e.binPath, args...)

	outBuff, errBuff := &bytes.Buffer{}, &bytes.Buffer{}
	cmd.Stdout = outBuff
	cmd.Stderr = errBuff

	err := cmd.Run()
	return outBuff.String(), errBuff.String(), err
}

func (e *KubectlEpinio) Create(resource string, args ...string) (string, string, error) {
	args = append([]string{"create", resource}, args...)
	return e.Run(args...)
}

func (e *KubectlEpinio) Get(resource string, args ...string) (string, string, error) {
	args = append([]string{"get", resource}, args...)
	return e.Run(args...)
}

func (e *KubectlEpinio) Delete(resource string, args ...string) (string, string, error) {
	args = append([]string{"delete", resource}, args...)
	return e.Run(args...)
}

func parseOutTable(out string) [][]string {
	outTable := [][]string{}

	out = strings.TrimSpace(out)
	rows := strings.Split(out, "\n")

	for _, row := range rows {
		rowCells := strings.FieldsFunc(row, func(r rune) bool {
			return r == '\t'
		})

		outTable = append(outTable, rowCells)

	}
	return outTable
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func setup() {
	_ = exec.Command(
		"kubectl", "delete", "secrets",
		"-n", "epinio",
		"-l", "epinio.io/api-user-credentials",
	).Run()

	_ = exec.Command(
		"kubectl", "delete", "configmaps",
		"-n", "epinio",
		"-l", "epinio.io/role",
	).Run()

	fmt.Printf("\033[1;33m%s\033[0m", "> Setup completed\n")
}

func teardown() {
	_ = exec.Command(
		"kubectl", "delete", "secrets",
		"-n", "epinio",
		"-l", "epinio.io/api-user-credentials",
	).Run()

	_ = exec.Command(
		"kubectl", "delete", "configmaps",
		"-n", "epinio",
		"-l", "epinio.io/role",
	).Run()

	fmt.Printf("\033[1;33m%s\033[0m", "> Teardown completed")
	fmt.Printf("\n")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
