package webhook

import (
	"testing"

	"github.com/kyma-project/kyma/components/function-controller/pkg/apis/serverless/v1alpha1"
	"github.com/kyma-project/kyma/components/function-controller/pkg/utils"
	"github.com/onsi/gomega"
)

func TestFunctionHandler_Default_Uninitialized_Values(t *testing.T) {
	type fields struct {
		Function    v1alpha1.Function
		RuntimeInfo utils.RuntimeInfo
	}

	rtInfo := utils.RuntimeInfo{
		Defaults: utils.DefaultConfig{
			Runtime:         "some",
			Size:            "random",
			TimeOut:         7312,
			FuncContentType: "values",
		},
	}

	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "should set uninitialized fields",
			fields: fields{
				Function:    v1alpha1.Function{},
				RuntimeInfo: rtInfo,
			},
		},
		{
			name: "should set fields with zeroed values",
			fields: fields{
				Function: v1alpha1.Function{
					Spec: v1alpha1.FunctionSpec{
						FunctionContentType: "",
						Size:                "",
						Runtime:             "",
						Timeout:             0,
					},
				},
				RuntimeInfo: rtInfo,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gomega.NewGomegaWithT(t)

			r := &FunctionHandler{
				Function:    tt.fields.Function,
				RuntimeInfo: tt.fields.RuntimeInfo,
			}

			g.Expect(r.Spec.Runtime).NotTo(gomega.Equal(r.RuntimeInfo.Defaults.Runtime))
			g.Expect(r.Spec.Size).NotTo(gomega.Equal(r.RuntimeInfo.Defaults.Size))
			g.Expect(r.Spec.Timeout).NotTo(gomega.Equal(r.RuntimeInfo.Defaults.TimeOut))
			g.Expect(r.Spec.FunctionContentType).NotTo(gomega.Equal(r.RuntimeInfo.Defaults.FuncContentType))

			r.Default()

			g.Expect(r.Spec.Runtime).To(gomega.Equal(r.RuntimeInfo.Defaults.Runtime))
			g.Expect(r.Spec.Size).To(gomega.Equal(r.RuntimeInfo.Defaults.Size))
			g.Expect(r.Spec.Timeout).To(gomega.Equal(r.RuntimeInfo.Defaults.TimeOut))
			g.Expect(r.Spec.FunctionContentType).To(gomega.Equal(r.RuntimeInfo.Defaults.FuncContentType))
		})
	}
}

func TestFunctionHandler_Default_Existing_Values(t *testing.T) {
	rtInfo := utils.RuntimeInfo{
		Defaults: utils.DefaultConfig{
			Runtime:         "some",
			Size:            "random",
			TimeOut:         7312,
			FuncContentType: "values",
		},
	}

	fn := v1alpha1.Function{
		Spec: v1alpha1.FunctionSpec{
			FunctionContentType: "existing",
			Size:                "real",
			Runtime:             "value",
			Timeout:             1000,
		},
	}

	t.Run("should not set already initialized fields", func(t *testing.T) {
		g := gomega.NewGomegaWithT(t)

		r := &FunctionHandler{
			Function:    fn,
			RuntimeInfo: rtInfo,
		}

		g.Expect(r.Spec.Runtime).NotTo(gomega.Equal(r.RuntimeInfo.Defaults.Runtime))
		g.Expect(r.Spec.Size).NotTo(gomega.Equal(r.RuntimeInfo.Defaults.Size))
		g.Expect(r.Spec.Timeout).NotTo(gomega.Equal(r.RuntimeInfo.Defaults.TimeOut))
		g.Expect(r.Spec.FunctionContentType).NotTo(gomega.Equal(r.RuntimeInfo.Defaults.FuncContentType))

		r.Default()

		g.Expect(r.Spec.Runtime).NotTo(gomega.Equal(r.RuntimeInfo.Defaults.Runtime))
		g.Expect(r.Spec.Size).NotTo(gomega.Equal(r.RuntimeInfo.Defaults.Size))
		g.Expect(r.Spec.Timeout).NotTo(gomega.Equal(r.RuntimeInfo.Defaults.TimeOut))
		g.Expect(r.Spec.FunctionContentType).NotTo(gomega.Equal(r.RuntimeInfo.Defaults.FuncContentType))
	})
}
