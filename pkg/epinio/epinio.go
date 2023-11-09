package epinio

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"

	"github.com/enrichman/kubectl-epinio/pkg/epinio/internal/names"
)

var ErrUserNotFound error = errors.New("user not found")

type KubeClient struct {
	kube kubernetes.Interface
}

func NewKubeClient(kubeClient kubernetes.Interface) (k *KubeClient) {
	return &KubeClient{kube: kubeClient}
}

func (k *KubeClient) GetUser(ctx context.Context, username string) (User, error) {
	users, err := k.ListUsers(ctx)
	if err != nil {
		return User{}, err
	}

	for _, u := range users {
		if u.Username == username {
			return u, nil
		}
	}

	return User{}, ErrUserNotFound
}

func (k *KubeClient) ListUsers(ctx context.Context) ([]User, error) {
	userSelector := labels.Set(map[string]string{
		"epinio.io/api-user-credentials": "true",
	}).AsSelector()

	secretClient := k.kube.CoreV1().Secrets("epinio")
	secretList, err := secretClient.List(ctx, metav1.ListOptions{LabelSelector: userSelector.String()})
	if err != nil {
		return nil, err
	}

	users := []User{}
	for _, sec := range secretList.Items {
		user := User{
			Username:          string(sec.Data["username"]),
			Password:          string(sec.Data["password"]),
			Role:              sec.Labels["epinio.io/role"],
			CreationTimestamp: sec.CreationTimestamp.Time,

			secret: sec.Name,
		}

		namespacesAll := string(sec.Data["namespaces"])
		namespacesAll = strings.TrimSpace(namespacesAll)
		if namespacesAll != "" {
			namespaces := strings.Split(namespacesAll, "\n")
			user.Namespaces = namespaces
		}

		rolesAll := strings.TrimSpace(sec.Annotations["epinio.io/roles"])
		if rolesAll != "" {
			roles := strings.Split(rolesAll, "\n")
			user.Namespaces = roles
		}

		users = append(users, user)
	}

	return users, nil
}

func (k *KubeClient) PatchUser(ctx context.Context, user User) error {
	patchSecretData := map[string][]byte{}

	if len(user.Namespaces) > 0 {
		nsData := strings.Join(user.Namespaces, "\n")
		patchSecretData["namespaces"] = []byte(nsData)
	}

	patch, err := json.Marshal(v1.Secret{Data: patchSecretData})
	if err != nil {
		return err
	}

	secretClient := k.kube.CoreV1().Secrets("epinio")
	_, err = secretClient.Patch(ctx, user.secret, types.StrategicMergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (k *KubeClient) CreateUser(ctx context.Context, user User) error {
	createSecretData := map[string][]byte{}

	createSecretData["username"] = []byte(user.Username)

	if user.Password != "" {
		createSecretData["password"] = []byte(user.Password)
	}

	if len(user.Namespaces) > 0 {
		nsData := strings.Join(user.Namespaces, "\n")
		createSecretData["namespaces"] = []byte(nsData)
	}

	userSecret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: names.GenerateResourceName("ruser", user.Username),
			Labels: map[string]string{
				"epinio.io/api-user-credentials": "true",
			},
		},
		Data: createSecretData,
	}

	secretClient := k.kube.CoreV1().Secrets("epinio")
	_, err := secretClient.Create(ctx, userSecret, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}
