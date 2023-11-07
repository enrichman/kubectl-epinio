package epinio

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

type KubeClient struct {
	kube kubernetes.Interface
}

func NewKubeClient(kubeClient kubernetes.Interface) (k *KubeClient) {
	return &KubeClient{kube: kubeClient}
}

func (k *KubeClient) ListUsers(ctx context.Context) ([]User, error) {
	userSelector := labels.Set(map[string]string{
		"epinio.io/api-user-credentials": "true",
	}).AsSelector()

	secretClient := k.kube.CoreV1().Secrets("epinio")
	secretList, err := secretClient.List(ctx, v1.ListOptions{LabelSelector: userSelector.String()})
	if err != nil {
		return nil, err
	}

	users := []User{}
	for _, sec := range secretList.Items {
		users = append(users, User{
			Username: string(sec.Data["username"]),
			Password: string(sec.Data["password"]),
		})
	}

	return users, nil
}
