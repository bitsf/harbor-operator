package portal

import (
	"context"

	"github.com/pkg/errors"

	containerregistryv1alpha2 "github.com/ovh/harbor-operator/api/v1alpha2"
)

func (r *Reconciler) AddResources(ctx context.Context, portal *containerregistryv1alpha2.Portal) error {
	service, err := r.GetService(ctx, portal)
	if err != nil {
		return errors.Wrap(err, "cannot get service")
	}

	_, err = r.Controller.AddBasicObjectToManage(ctx, service)
	if err != nil {
		return errors.Wrapf(err, "cannot add service %+v", service)
	}

	deployment, err := r.GetDeployment(ctx, portal)
	if err != nil {
		return errors.Wrap(err, "cannot get deployment")
	}

	_, err = r.Controller.AddBasicObjectToManage(ctx, deployment)
	if err != nil {
		return errors.Wrapf(err, "cannot add deployment %+v", deployment)
	}

	return errors.New("not yet implemented")
}
