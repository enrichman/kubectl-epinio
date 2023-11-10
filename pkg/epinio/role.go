package epinio

import (
	"context"
	"strconv"
	"strings"

	"github.com/enrichman/kubectl-epinio/pkg/epinio/internal/names"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

var Actions = []string{
	"namespace",
	"namespace_read",
	"namespace_write",
	"app",
	"app_read",
	"app_write",
	"app_logs",
	"app_exec",
	"app_portforward",
	"configuration",
	"configuration_read",
	"configuration_write",
	"service",
	"service_read",
	"service_write",
	"service_portforward",
	"gitconfig",
	"gitconfig_read",
	"gitconfig_write",
	"export_registries_read",
}

type Role struct {
	ID      string
	Default bool
	Name    string
	Actions []string
}

func (r Role) GetID() string {
	return r.ID
}

func (k *KubeClient) ListRoles(ctx context.Context) ([]Role, error) {
	roleSelector := labels.Set(map[string]string{
		"epinio.io/role": "true",
	}).AsSelector()

	opts := metav1.ListOptions{LabelSelector: roleSelector.String()}
	roleList, err := k.kube.CoreV1().ConfigMaps("epinio").List(ctx, opts)
	if err != nil {
		return nil, err
	}

	roles := []Role{{ID: "admin", Name: "Admin Role"}}
	for _, cm := range roleList.Items {
		actions := []string{}

		if cm.Data["actions"] != "" {
			splitted := strings.Split(cm.Data["actions"], "\n")
			actions = append(actions, splitted...)
		}

		roles = append(roles, Role{
			ID:      cm.Data["id"],
			Name:    cm.Data["name"],
			Actions: actions,
		})
	}

	return roles, nil
}

func (k *KubeClient) CreateRole(ctx context.Context, role Role) error {
	createRoleData := map[string]string{}

	createRoleData["id"] = role.ID

	if role.Default {
		createRoleData["default"] = strconv.FormatBool(role.Default)
	}

	if len(role.Actions) > 0 {
		actionsData := strings.Join(role.Actions, "\n")
		createRoleData["actions"] = actionsData
	}

	roleConfigMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: names.GenerateResourceName("epinio", role.ID, "role"),
			Labels: map[string]string{
				"epinio.io/role": "true",
			},
		},
		Data: createRoleData,
	}

	cmClient := k.kube.CoreV1().ConfigMaps("epinio")
	_, err := cmClient.Create(ctx, roleConfigMap, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// DeleteRole deletes a role by id.
func (k *KubeClient) DeleteRole(ctx context.Context, id string) error {
	cmClient := k.kube.CoreV1().ConfigMaps("epinio")

	name := names.GenerateResourceName("epinio", id, "role")

	return cmClient.Delete(ctx, name, metav1.DeleteOptions{})
}
