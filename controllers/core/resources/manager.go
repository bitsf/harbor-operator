package notaryserverresources

import (
	containerregistryv1alpha2 "github.com/ovh/harbor-operator/api/v1alpha2"
)

type Manager struct {
	Core *containerregistryv1alpha2.Core
}
