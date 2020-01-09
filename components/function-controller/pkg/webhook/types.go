package webhook

import (
	"github.com/kyma-project/kyma/components/function-controller/pkg/apis/serverless/v1alpha1"
	runtimeUtil "github.com/kyma-project/kyma/components/function-controller/pkg/utils"
)

type FunctionHandler struct {
	v1alpha1.Function
	// Client      client.Client
	RuntimeInfo runtimeUtil.RuntimeInfo
}
