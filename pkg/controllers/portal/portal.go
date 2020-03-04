package portal

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ovh/harbor-operator/controllers/portal"
	"github.com/ovh/harbor-operator/pkg/controllers"
	"github.com/ovh/harbor-operator/pkg/controllers/config"
)

const (
	Name         = "portal"
	ConfigPrefix = Name + "-controller"
)

func New(ctx context.Context, version string) (controllers.Controller, error) {
	config, err := config.GetConfig(ConfigPrefix)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get configuration")
	}

	return portal.New(ctx, Name, version, config)
}
