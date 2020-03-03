package portal

import (
	"context"

	"github.com/pkg/errors"

	containerregistryv1alpha2 "github.com/ovh/harbor-operator/api/v1alpha2"
)

func (r *Reconciler) AddResources(ctx context.Context, portal *containerregistryv1alpha2.Portal) error {
	service, err := r.GetService(ctx, portal)
	if err != nil {
		return errors.Wrap(err, "cannot get configMap")
	}

	err = r.Controller.AddResourceToManage(ctx, service)
	if err != nil {
		return errors.Wrapf(err, "cannot add resource %+v", cm)
	}

	return errors.New("not yet implemented")
}
