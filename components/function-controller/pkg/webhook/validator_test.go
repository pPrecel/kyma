package webhook

import (
	"testing"

	"github.com/kyma-project/kyma/components/function-controller/pkg/apis/serverless/v1alpha1"
	"github.com/kyma-project/kyma/components/function-controller/pkg/utils"
)

func TestFunctionHandler_ValidateFunctionFormat(t *testing.T) {
	type fields struct {
		Function    v1alpha1.Function
		RuntimeInfo utils.RuntimeInfo
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "tada",
			wantErr: true,
			fields: fields{
				Function: v1alpha1.Function{
					Spec: v1alpha1.FunctionSpec{
						Size: "size3",
					},
				},
				RuntimeInfo: utils.RuntimeInfo{
					FuncSizes: []utils.FuncSize{
						{Size: "size1"},
						{Size: "size2"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &FunctionHandler{
				Function:    tt.fields.Function,
				RuntimeInfo: tt.fields.RuntimeInfo,
			}
			if err := r.ValidateFunctionFormat(); (err != nil) != tt.wantErr {
				t.Errorf("ValidateFunctionFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
