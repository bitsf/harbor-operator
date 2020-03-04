package jobservice

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ovh/harbor-operator/controllers/jobservice"
	"github.com/ovh/harbor-operator/pkg/controllers"
	"github.com/ovh/harbor-operator/pkg/controllers/config"
)

const (
	Name         = "jobservice"
	ConfigPrefix = Name + "-controller"
)

func New(ctx context.Context, version string) (controllers.Controller, error) {
	config, err := config.GetConfig(ConfigPrefix)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get configuration")
	}

	return jobservice.New(ctx, Name, version, config)
}
