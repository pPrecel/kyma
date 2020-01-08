package webhook

import (
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var _ webhook.Defaulter = &FunctionHandler{}

// +kubebuilder:webhook:path=/function-webhook.serverless.kyma-project.io,mutating=true,failurePolicy=fail,groups=serverless.kyma-project.io,resources=functions,verbs=create;update,versions=v1alpha1,name=mfunction.kb.io

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *FunctionHandler) Default() {
	if r.Spec.Size == "" {
		r.Spec.Size = r.RuntimeInfo.Defaults.Size
	}
	if r.Spec.Runtime == "" {
		r.Spec.Runtime = r.RuntimeInfo.Defaults.Runtime
	}
	if r.Spec.Timeout == 0 {
		r.Spec.Timeout = r.RuntimeInfo.Defaults.TimeOut
	}
	if r.Spec.FunctionContentType == "" {
		r.Spec.FunctionContentType = r.RuntimeInfo.Defaults.FuncContentType
	}
}