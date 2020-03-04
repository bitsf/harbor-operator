package registry

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ovh/harbor-operator/controllers/registry"
	"github.com/ovh/harbor-operator/pkg/controllers"
	"github.com/ovh/harbor-operator/pkg/controllers/config"
)

const (
	Name         = "registry"
	ConfigPrefix = Name + "-controller"
)

func New(ctx context.Context, version string) (controllers.Controller, error) {
	config, err := config.GetConfig(ConfigPrefix)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get configuration")
	}

	return registry.New(ctx, Name, version, config)
}
