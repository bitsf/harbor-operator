package harbor

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	containerregistryv1alpha2 "github.com/ovh/harbor-operator/api/v1alpha2"
	"github.com/ovh/harbor-operator/pkg/factories/application"
	"github.com/ovh/harbor-operator/pkg/factories/logger"
)

// +kubebuilder:rbac:groups=containerregistry.ovhcloud.com,resources=harbors,verbs=get;list;watch
// +kubebuilder:rbac:groups=containerregistry.ovhcloud.com,resources=harbors/status,verbs=get;update;patch

func (r *Reconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.TODO()
	application.SetName(&ctx, r.GetName())
	application.SetVersion(&ctx, r.GetVersion())

	span, ctx := opentracing.StartSpanFromContext(ctx, "reconcile", opentracing.Tags{
		"Harbor.Namespace": req.Namespace,
		"Harbor.Name":      req.Name,
	})
	defer span.Finish()

	span.LogFields(
		log.String("Harbor.Namespace", req.Namespace),
		log.String("Harbor.Name", req.Name),
	)

	reqLogger := r.Log.WithValues("Request", req.NamespacedName, "Harbor.Namespace", req.Namespace, "Harbor.Name", req.Name)

	logger.Set(&ctx, reqLogger)

	// Fetch the Harbor instance
	harbor := &containerregistryv1alpha2.Harbor{}

	ok, err := r.Controller.GetAndFilter(ctx, req.NamespacedName, harbor)
	if err != nil {
		// Error reading the object
		return reconcile.Result{}, err
	}
	if !ok {
		// Request object not found, could have been deleted after reconcile request.
		// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
		reqLogger.Info("Harbor does not exists")
		return reconcile.Result{}, nil
	}

	err = r.AddResources(ctx, harbor)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "cannot add resources")
	}

	return r.Controller.Reconcile(ctx, harbor)
}
