package epinio

import (
	"context"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type Role struct {
	ID      string
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

	roles := []Role{}
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
