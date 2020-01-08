package webhook

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

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
	for _, functionSize := range r.RuntimeInfo.FuncSizes {
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
	for _, runtime := range r.RuntimeInfo.AvailableRuntimes {
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
	for _, functionContentType := range r.RuntimeInfo.FuncTypes {
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
