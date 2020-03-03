package harborresources

import (
	containerregistryv1alpha2 "github.com/ovh/harbor-operator/api/v1alpha2"
)

type Manager struct {
	Harbor *containerregistryv1alpha2.Harbor
}
