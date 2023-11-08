package epinio

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (k *KubeClient) ListNamespaces(ctx context.Context) ([]string, error) {
	namespaceSelector := labels.Set(map[string]string{
		"app.kubernetes.io/component": "epinio-namespace",
	}).AsSelector()

	opts := metav1.ListOptions{LabelSelector: namespaceSelector.String()}

	namespaceList, err := k.kube.CoreV1().Namespaces().List(ctx, opts)
	if err != nil {
		return nil, err
	}

	namespaces := []string{}
	for _, ns := range namespaceList.Items {
		namespaces = append(namespaces, ns.Name)
	}

	return namespaces, nil
}
