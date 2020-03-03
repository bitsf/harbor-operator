package portal

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ovh/harbor-operator/controllers/harbor"
	"github.com/ovh/harbor-operator/pkg/controllers/config"
)

const (
	ConfigPrefix = "portal-controller"
)

func New(ctx context.Context, name, version string) (*harbor.Reconciler, error) {
	config, err := config.GetConfig(ConfigPrefix)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get configuration")
	}

	return harbor.New(ctx, name, version, config)
}
