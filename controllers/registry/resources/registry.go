package registryresources

import (
	containerregistryv1alpha2 "github.com/ovh/harbor-operator/api/v1alpha2"
)

type Manager struct {
	Registry *containerregistryv1alpha2.Registry
}
