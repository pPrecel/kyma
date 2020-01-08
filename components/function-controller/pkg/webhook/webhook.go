package webhook

import (
	"context"
	"fmt"
	"os"

	runtimeUtil "github.com/kyma-project/kyma/components/function-controller/pkg/utils"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

// log is for logging in this package.
var webhookLog = logf.Log.WithName("webhook")

var (
	// name of function config
	fnConfigName = getEnvDefault("CONTROLLER_CONFIGMAP", "fn-config")

	// namespace of function config
	fnConfigNamespace = getEnvDefault("CONTROLLER_CONFIGMAP_NS", "default")
)

func getEnvDefault(envName string, defaultValue string) string {
	// use default value if environment variable is empty
	var value string
	if value = os.Getenv(envName); value == "" {
		return defaultValue
	}
	return value
}

func (r *FunctionHandler) SetupWebhookWithManager(mgr ctrl.Manager) error {
	var err error
	rntInfo, err := getRuntimeInfo(r.Client, fnConfigName, fnConfigNamespace)
	if err != nil {
		return err
	}

	r.RuntimeInfo = *rntInfo

	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

func getRuntimeInfo(cl client.Client, fnConfigName, fnConfigNamespace string) (*runtimeUtil.RuntimeInfo, error) {
	cm := &corev1.ConfigMap{}

	err := cl.Get(context.TODO(),
		client.ObjectKey{
			Name:      fnConfigName,
			Namespace: fnConfigNamespace,
		},
		cm,
	)

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("while reading ConfigMap %s from namespace %s", fnConfigName, fnConfigNamespace))
	}

	rnInfo, err := runtimeUtil.New(cm)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("while creating RuntimeInfo from ConfigMap %s in namespace %s", cm.Name, cm.Namespace))
	}
	return rnInfo, nil
}
