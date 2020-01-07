/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mutating

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	serverlessv1alpha1 "github.com/kyma-project/kyma/components/function-controller/api/v1alpha1"
	runtimeUtil "github.com/kyma-project/kyma/components/function-controller/utils"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	log = logf.Log.WithName("webhook")
)

const WebhookName = "mutating-create-function"

func init() {
	log.Info("init")
	if HandlerMap[WebhookName] == nil {
		HandlerMap[WebhookName] = []admission.Handler{}
	}
	HandlerMap[WebhookName] = append(HandlerMap[WebhookName], &FunctionCreateHandler{})
}

// FunctionCreateHandler handles Function
type FunctionCreateHandler struct {
	Client client.Client

	// Decoder decodes objects
	Decoder admission.Decoder
}

// Mutates function values
func (h *FunctionCreateHandler) mutatingFunctionFn(obj *serverlessv1alpha1.Function, rnInfo *runtimeUtil.RuntimeInfo) {
	if obj.Spec.Size == "" {
		obj.Spec.Size = rnInfo.Defaults.Size
	}
	if obj.Spec.Runtime == "" {
		obj.Spec.Runtime = rnInfo.Defaults.Runtime
	}
	if obj.Spec.Timeout == 0 {
		obj.Spec.Timeout = rnInfo.Defaults.TimeOut
	}
	if obj.Spec.FunctionContentType == "" {
		obj.Spec.FunctionContentType = rnInfo.Defaults.FuncContentType
	}
}

// Validate function values and return an error if the function is not valid
func (h *FunctionCreateHandler) validateFunctionFn(obj *serverlessv1alpha1.Function, rnInfo *runtimeUtil.RuntimeInfo) error {
	// function size
	var isValidFunctionSize bool
	var functionSizes []string
	for _, functionSize := range rnInfo.FuncSizes {
		functionSizes = append(functionSizes, functionSize.Size)
		if obj.Spec.Size == functionSize.Size {
			isValidFunctionSize = true
			break
		}
	}

	if !isValidFunctionSize {
		return fmt.Errorf("size should be one of %q (got %q)",
			functionSizes, obj.Spec.Size)
	}

	// function serverless
	var isValidRuntime bool
	var runtimes []string
	for _, runtime := range rnInfo.AvailableRuntimes {
		runtimes = append(runtimes, runtime.ID)
		if obj.Spec.Runtime == runtime.ID {
			isValidRuntime = true
			break
		}
	}

	if !isValidRuntime {
		return fmt.Errorf("runtime should be one of %q (got %q)",
			runtimes, obj.Spec.Runtime)
	}

	// function content type
	var isValidFunctionContentType bool
	var functionContentTypes []string
	for _, functionContentType := range rnInfo.FuncTypes {
		functionContentTypes = append(functionContentTypes, functionContentType.Type)
		if obj.Spec.FunctionContentType == functionContentType.Type {
			isValidFunctionContentType = true
			break
		}
	}

	if !isValidFunctionContentType {
		return fmt.Errorf("functionContentType should be one of %q (got %q)",
			functionContentTypes, obj.Spec.FunctionContentType)
	}

	return nil
}

var _ admission.Handler = &FunctionCreateHandler{}

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

// getRuntimeConfig returns the Function Controller ConfigMap from the cluster.
func (h *FunctionCreateHandler) getRuntimeConfig() (*corev1.ConfigMap, error) {
	cm := &corev1.ConfigMap{}

	err := h.Client.Get(context.TODO(),
		client.ObjectKey{
			Name:      fnConfigName,
			Namespace: fnConfigNamespace,
		},
		cm,
	)

	if err != nil {
		return nil, err
	}
	return cm, nil
}

// Handle handles admission requests.
func (h *FunctionCreateHandler) Handle(ctx context.Context, req admission.Request) admission.Response {
	log.Info("received admission request", "request", req)

	fnConfig, err := h.getRuntimeConfig()
	if err != nil {
		log.Error(err, "Error reading controller configuration", "namespace", fnConfigNamespace, "name", fnConfigName)
		return admission.Errored(http.StatusInternalServerError, err)
	}

	rnInfo, err := runtimeUtil.New(fnConfig)
	if err != nil {
		log.Error(err, "Error creating RuntimeInfo", "namespace", fnConfig.Namespace, "name", fnConfig.Name)
		return admission.Errored(http.StatusBadRequest, err)
	}
	obj := &serverlessv1alpha1.Function{}

	if err := h.Decoder.Decode(req, obj); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}
	copy := obj.DeepCopy()

	// mutate values
	h.mutatingFunctionFn(copy, rnInfo)

	// validate function and return an error describing the validation error if validation fails
	if err := h.validateFunctionFn(copy, rnInfo); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	marshalledObj, err := json.Marshal(copy)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}
	return admission.PatchResponseFromRaw(req.Object.Raw, marshalledObj)
}

var _ inject.Client = &FunctionCreateHandler{}

// InjectClient injects the client into the FunctionCreateHandler
func (h *FunctionCreateHandler) InjectClient(c client.Client) error {
	h.Client = c
	return nil
}

// InjectDecoder injects the decoder into the FunctionCreateHandler
func (h *FunctionCreateHandler) InjectDecoder(d admission.Decoder) error {
	h.Decoder = d
	return nil
}
