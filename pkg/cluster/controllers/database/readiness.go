package database

import (
	"context"
	"errors"
	"fmt"

	goharborv1alpha2 "github.com/goharbor/harbor-operator/apis/goharbor.io/v1alpha2"
	harbormetav1 "github.com/goharbor/harbor-operator/apis/meta/v1alpha1"
	"github.com/goharbor/harbor-operator/pkg/cluster/controllers/database/api"
	"github.com/goharbor/harbor-operator/pkg/cluster/lcm"
	corev1 "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	labels1 "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	HarborCore         = "core"
	HarborClair        = "clair"
	HarborNotaryServer = "notaryServer"
	HarborNotarySigner = "notarySigner"

	CoreDatabase         = "core"
	ClairDatabase        = "clair"
	NotaryServerDatabase = "notaryserver"
	NotarySignerDatabase = "notarysigner"
	DefaultDatabaseUser  = "harbor"

	CoreSecretName         = "core"
	ClairSecretName        = "clair"
	NotaryServerSecretName = "notary-server"
	NotarySignerSecretName = "notary-signer"
)

// Readiness reconcile will check postgre sql cluster if that has available.
// It does:
// - create postgre connection pool
// - ping postgre server
// - return postgre properties if postgre has available.
func (p *PostgreSQLController) Readiness(ctx context.Context) (*lcm.CRStatus, error) {
	var (
		conn *Connect
		err  error
	)

	name := p.HarborCluster.Name

	conn, err = p.GetInClusterDatabaseInfo()
	if err != nil {
		return nil, err
	}

	var pg api.Postgresql
	if err := runtime.DefaultUnstructuredConverter.
		FromUnstructured(p.ActualCR.UnstructuredContent(), &pg); err != nil {
		return nil, err
	}

	if pg.Status.PostgresClusterStatus != "Running" {
		return nil, errors.New("database is not ready")
	}

	p.Log.Info("Database already ready.",
		"namespace", p.HarborCluster.Namespace,
		"name", p.HarborCluster.Name)

	properties := &lcm.Properties{}

	if err := p.DeployComponentSecret(conn, getDatabasePasswordRefName(name)); err != nil {
		return nil, err
	}

	addProperties(name, conn, properties)

	crStatus := lcm.New(goharborv1alpha2.DatabaseReady).
		WithStatus(corev1.ConditionTrue).
		WithReason("database already ready").
		WithMessage("harbor component database secrets are already create.").
		WithProperties(*properties)

	return crStatus, nil
}

func addProperties(name string, conn *Connect, properties *lcm.Properties) {
	db := getHarborDatabaseSpec(name, conn)
	properties.Add(lcm.DatabasePropertyName, db)
}

func getHarborDatabaseSpec(name string, conn *Connect) *goharborv1alpha2.HarborDatabaseSpec {
	return &goharborv1alpha2.HarborDatabaseSpec{
		PostgresCredentials: harbormetav1.PostgresCredentials{
			Username:    DefaultDatabaseUser,
			PasswordRef: getDatabasePasswordRefName(name),
		},
		Hosts: []harbormetav1.PostgresHostSpec{
			{
				Host: conn.Host,
				Port: InClusterDatabasePortInt32,
			},
		},
		SSLMode: harbormetav1.PostgresSSLModeDisable,
	}
}

func getDatabasePasswordRefName(name string) string {
	return fmt.Sprintf("%s-%s-%s", name, "database", "password")
}

// DeployComponentSecret deploy harbor component database secret.
func (p *PostgreSQLController) DeployComponentSecret(conn *Connect, secretName string) error {
	secret := &corev1.Secret{}
	sc := p.GetDatabaseSecret(conn, secretName)

	if err := controllerutil.SetControllerReference(p.HarborCluster, sc, p.Scheme); err != nil {
		return err
	}

	err := p.Client.Get(types.NamespacedName{Name: secretName, Namespace: p.HarborCluster.Namespace}, secret)
	if kerr.IsNotFound(err) {
		p.Log.Info("Creating Harbor Component Secret",
			"namespace", p.HarborCluster.Namespace,
			"name", secretName)

		return p.Client.Create(sc)
	} else if err != nil {
		return err
	}

	return nil
}

// GetInClusterDatabaseInfo returns inCluster database connection client.
func (p *PostgreSQLController) GetInClusterDatabaseInfo() (*Connect, error) {
	var (
		connect *Connect
		err     error
	)

	pw, err := p.GetInClusterDatabasePassword()
	if err != nil {
		return connect, err
	}

	if connect, err = p.GetInClusterDatabaseConn(p.resourceName(), pw); err != nil {
		return connect, err
	}

	return connect, nil
}

// GetInClusterDatabaseConn returns inCluster database connection info.
func (p *PostgreSQLController) GetInClusterDatabaseConn(name, pw string) (*Connect, error) {
	host, err := p.GetInClusterHost(name)
	if err != nil {
		return nil, err
	}

	conn := &Connect{
		Host:     host,
		Port:     InClusterDatabasePort,
		Password: pw,
		Username: DefaultDatabaseUser,
		Database: CoreDatabase,
	}

	return conn, nil
}

func GenInClusterPasswordSecretName(user, crName string) string {
	return fmt.Sprintf("%s.%s.credentials", user, crName)
}

// GetInClusterHost returns the Database master pod ip or service name.
func (p *PostgreSQLController) GetInClusterHost(name string) (string, error) {
	var (
		url string
		err error
	)

	_, err = rest.InClusterConfig()
	if err != nil {
		url, err = p.GetMasterPodsIP()
		if err != nil {
			return url, err
		}
	} else {
		url = fmt.Sprintf("%s.%s.svc", name, p.HarborCluster.Namespace)
	}

	return url, nil
}

// GetInClusterDatabasePassword is get inCluster postgresql password.
func (p *PostgreSQLController) GetInClusterDatabasePassword() (string, error) {
	var pw string

	secretName := GenInClusterPasswordSecretName(DefaultDatabaseUser, p.resourceName())

	secret, err := p.GetSecret(secretName)
	if err != nil {
		return pw, err
	}

	for k, v := range secret {
		if k == InClusterDatabasePasswordKey {
			pw = string(v)

			return pw, nil
		}
	}

	return pw, nil
}

// GetStatefulSetPods returns the postgresql master pod.
func (p *PostgreSQLController) GetStatefulSetPods() (*corev1.PodList, error) {
	label := map[string]string{
		"application":  "spilo",
		"cluster-name": p.resourceName(),
		"spilo-role":   "master",
	}

	opts := &client.ListOptions{}
	set := labels1.SelectorFromSet(label)
	opts.LabelSelector = set
	pod := &corev1.PodList{}

	if err := p.Client.List(opts, pod); err != nil {
		p.Log.Error(err, "fail to get pod.",
			"namespace", p.HarborCluster.Namespace, "name", p.resourceName())

		return nil, err
	}

	return pod, nil
}

// GetMasterPodsIP returns postgresql master node ip.
func (p *PostgreSQLController) GetMasterPodsIP() (string, error) {
	var masterIP string

	podList, err := p.GetStatefulSetPods()
	if err != nil {
		return masterIP, err
	}

	if len(podList.Items) > 1 {
		return masterIP, errors.New("the number of master node copies cannot exceed 1")
	}

	for _, p := range podList.Items {
		if p.DeletionTimestamp != nil {
			continue
		}

		masterIP = p.Status.PodIP
	}

	return masterIP, nil
}
