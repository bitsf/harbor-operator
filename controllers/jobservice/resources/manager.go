package jobserviceresources

import (
	containerregistryv1alpha2 "github.com/ovh/harbor-operator/api/v1alpha2"
)

type Manager struct {
	JobService *containerregistryv1alpha2.JobService
}
