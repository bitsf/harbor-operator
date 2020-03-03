package registryctlresources

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/sethvargo/go-password/password"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	containerregistryv1alpha2 "github.com/ovh/harbor-operator/api/v1alpha2"
	"github.com/ovh/harbor-operator/pkg/factories/application"
)

const (
	keyLength = 15
)

func (m *Manager) GetSecrets(ctx context.Context) ([]*corev1.Secret, error) {
	operatorName := application.GetName(ctx)
	harborName := m.RegistryController.Name

	return []*corev1.Secret{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      m.RegistryController.Name,
				Namespace: m.RegistryController.Namespace,
				Labels: map[string]string{
					"app":      containerregistryv1alpha2.RegistryName,
					"harbor":   harborName,
					"operator": operatorName,
				},
			},
			Type: corev1.SecretTypeOpaque,
			StringData: map[string]string{
				"REGISTRY_HTTP_SECRET": password.MustGenerate(keyLength, 5, 5, false, true),
			},
		},
	}, nil
}

func (m *Manager) GetSecretsCheckSum() string {
	// TODO get generation of the secrets
	value := fmt.Sprintf("%s\n%s", m.RegistryController.Spec.CacheSecret, m.RegistryController.Spec.StorageSecret)
	sum := sha256.New().Sum([]byte(value))

	return fmt.Sprintf("%x", sum)
}
