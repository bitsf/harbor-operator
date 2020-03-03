package registryctlresources

import (
	containerregistryv1alpha2 "github.com/ovh/harbor-operator/api/v1alpha2"
)

type Manager struct {
	RegistryController *containerregistryv1alpha2.RegistryController
}
