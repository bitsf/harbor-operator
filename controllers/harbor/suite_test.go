package harbor

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	// +kubebuilder:scaffold:imports

	"github.com/goharbor/harbor-operator/pkg/factories/logger"
	"github.com/goharbor/harbor-operator/pkg/scheme"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t, "HarborController", []Reporter{envtest.NewlineReporter{}})
}

func setupTest(ctx context.Context) (*Reconciler, context.Context) {
	logger.Set(&ctx, zap.LoggerTo(GinkgoWriter, true))

	s, err := scheme.New(ctx)
	Expect(err).ToNot(HaveOccurred(), "failed to initialize scheme")

	return &Reconciler{
		Scheme: s,
	}, ctx
}
