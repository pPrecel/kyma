package webhook

import (
	"context"
	"fmt"
	"os"

	runtimeUtil "github.com/kyma-project/kyma/components/function-controller/pkg/utils"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var webhookLog = logf.Log.WithName("webhook")

var rntInfo runtimeUtil.RuntimeInfo

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

func (r FunctionHandler) SetupWebhookWithManager(mgr ctrl.Manager) error {
	var err error
	info, err := getRuntimeInfo(mgr.GetClient(), fnConfigName, fnConfigNamespace)
	if err != nil {
		return err
	}

	// TODO do not use global variable
	rntInfo = *info

	return ctrl.NewWebhookManagedBy(mgr).
		For(&r).
		Complete()
}

var _ webhook.Defaulter = &FunctionHandler{}

// +kubebuilder:webhook:path=/function-webhook.serverless.kyma-project.io,mutating=true,failurePolicy=fail,groups=serverless.kyma-project.io,resources=functions,verbs=create;update,versions=v1alpha1,name=mfunction.kb.io

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *FunctionHandler) Default() {
	if r.Spec.Size == "" {
		r.Spec.Size = rntInfo.Defaults.Size
	}
	if r.Spec.Runtime == "" {
		r.Spec.Runtime = rntInfo.Defaults.Runtime
	}
	if r.Spec.Timeout == 0 {
		r.Spec.Timeout = rntInfo.Defaults.TimeOut
	}
	if r.Spec.FunctionContentType == "" {
		r.Spec.FunctionContentType = rntInfo.Defaults.FuncContentType
	}
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-serverless-kyma-project-io-v1alpha1-function,mutating=false,failurePolicy=fail,groups=serverless.kyma-project.io,resources=functions,versions=v1alpha1,name=vfunction.kb.io

var _ webhook.Validator = &FunctionHandler{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *FunctionHandler) ValidateCreate() error {
	webhookLog.Info("validate create", "name", r.Name)
	return r.ValidateFunctionFormat()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *FunctionHandler) ValidateUpdate(old runtime.Object) error {
	webhookLog.Info("validate update", "name", r.Name)
	return r.ValidateFunctionFormat()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *FunctionHandler) ValidateDelete() error {
	return nil
}

func (r *FunctionHandler) ValidateFunctionFormat() error {
	// function size
	var isValidFunctionSize bool
	var functionSizes []string
	for _, functionSize := range rntInfo.FuncSizes {
		functionSizes = append(functionSizes, functionSize.Size)
		if r.Spec.Size == functionSize.Size {
			isValidFunctionSize = true
			break
		}
	}

	if !isValidFunctionSize {
		return fmt.Errorf("size should be one of %q (got %q)",
			functionSizes, r.Spec.Size)
	}

	// function serverless
	var isValidRuntime bool
	var runtimes []string
	for _, runtime := range rntInfo.AvailableRuntimes {
		runtimes = append(runtimes, runtime.ID)
		if r.Spec.Runtime == runtime.ID {
			isValidRuntime = true
			break
		}
	}

	if !isValidRuntime {
		return fmt.Errorf("runtime should be one of %q (got %q)",
			runtimes, r.Spec.Runtime)
	}

	// function content type
	var isValidFunctionContentType bool
	var functionContentTypes []string
	for _, functionContentType := range rntInfo.FuncTypes {
		functionContentTypes = append(functionContentTypes, functionContentType.Type)
		if r.Spec.FunctionContentType == functionContentType.Type {
			isValidFunctionContentType = true
			break
		}
	}

	if !isValidFunctionContentType {
		return fmt.Errorf("functionContentType should be one of %q (got %q)",
			functionContentTypes, r.Spec.FunctionContentType)
	}

	return nil
}

func getRuntimeInfo(cl client.Client, fnConfigName, fnConfigNamespace string) (*runtimeUtil.RuntimeInfo, error) {
	// I don't like this, I'd rather have this somehow injected like it used to be
	// but I don't know whether we
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
