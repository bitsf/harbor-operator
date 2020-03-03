package chartmuseumresources

import (
	containerregistryv1alpha2 "github.com/ovh/harbor-operator/api/v1alpha2"
)

type Manager struct {
	ChartMuseum *containerregistryv1alpha2.ChartMuseum
}
