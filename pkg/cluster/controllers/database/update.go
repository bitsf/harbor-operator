package database

import (
	"fmt"

	"github.com/goharbor/harbor-operator/pkg/cluster/controllers/database/api"
	"github.com/goharbor/harbor-operator/pkg/cluster/lcm"
	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Update reconcile will update PostgreSQL CR.
func (p *PostgreSQLController) Update() (*lcm.CRStatus, error) {
	name := fmt.Sprintf("%s-%s", p.HarborCluster.Namespace, p.HarborCluster.Name)

	crdClient := p.DClient.WithResource(databaseGVR).WithNamespace(p.HarborCluster.Namespace)

	if p.ExpectCR == nil {
		return databaseUnknownStatus(), nil
	}

	var actualCR, expectCR api.Postgresql

	if err := runtime.DefaultUnstructuredConverter.
		FromUnstructured(p.ActualCR.UnstructuredContent(), &actualCR); err != nil {
		return databaseNotReadyStatus(DefaultUnstructuredConverterError, err.Error()), err
	}

	if err := runtime.DefaultUnstructuredConverter.
		FromUnstructured(p.ExpectCR.UnstructuredContent(), &expectCR); err != nil {
		return databaseNotReadyStatus(DefaultUnstructuredConverterError, err.Error()), err
	}

	if !IsEqual(expectCR, actualCR) {
		p.Log.Info(
			"Update Database resource",
			"namespace", p.HarborCluster.Namespace, "name", name,
		)

		expectCR.ObjectMeta.SetResourceVersion(actualCR.ObjectMeta.GetResourceVersion())

		data, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&expectCR)
		if err != nil {
			return databaseNotReadyStatus(DefaultUnstructuredConverterError, err.Error()), err
		}

		_, err = crdClient.Update(&unstructured.Unstructured{Object: data}, metav1.UpdateOptions{})
		if err != nil {
			return databaseNotReadyStatus(UpdateDatabaseCrError, err.Error()), err
		}
	}

	return databaseUnknownStatus(), nil
}

// isEqual check whether cache cr is equal expect.
func IsEqual(actualCR, expectCR api.Postgresql) bool {
	return cmp.Equal(expectCR.DeepCopy().Spec, actualCR.DeepCopy().Spec)
}
