package registryctl

import (
	"context"

	"github.com/pkg/errors"

	containerregistryv1alpha2 "github.com/ovh/harbor-operator/api/v1alpha2"
)

func (r *Reconciler) InitResources() error {
	return errors.Wrap(r.InitConfigMaps(), "configmaps")
}

func (r *Reconciler) AddResources(ctx context.Context, registryctl *containerregistryv1alpha2.RegistryController) error {
	cm, err := r.GetConfigMap(ctx, registryctl)
	if err != nil {
		return errors.Wrap(err, "cannot get configMap")
	}

	_, err = r.Controller.AddInstantResourceToManage(ctx, cm)
	if err != nil {
		return errors.Wrapf(err, "cannot add configMap %+v", cm)
	}

	return errors.New("not yet implemented")
}
