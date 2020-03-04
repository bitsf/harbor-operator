package registry

import (
	"context"
	"fmt"

	"github.com/sethvargo/go-password/password"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	containerregistryv1alpha2 "github.com/ovh/harbor-operator/api/v1alpha2"
)

const (
	keyLength   = 15
	numDigits   = 5
	numSymbols  = 3
	noUpper     = false
	allowRepeat = true
)

func (r *Reconciler) GetSecret(ctx context.Context, registry *containerregistryv1alpha2.Registry) (*corev1.Secret, error) {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-registry", registry.GetName()),
			Namespace: registry.GetNamespace(),
		},
		Type: corev1.SecretTypeOpaque,
		StringData: map[string]string{
			"REGISTRY_HTTP_SECRET": password.MustGenerate(keyLength, numDigits, numSymbols, noUpper, allowRepeat),
		},
	}, nil
}
