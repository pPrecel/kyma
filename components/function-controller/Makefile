PKG = github.com/kyma-project/kyma/components/function-controller

# Image URL to use all building/pushing image targets
APP_NAME = function-controller
APP_PATH = components/$(APP_NAME)
BUILDPACK = eu.gcr.io/kyma-project/test-infra/buildpack-golang-kubebuilder:v20190208-813daef
SCRIPTS_DIR = $(realpath $(shell pwd)/../..)/common/makefiles


override ENTRYPOINT := ./main.go

include $(SCRIPTS_DIR)/generic-make-go.mk


# Image URL to use all building/pushing image targets
IMG ?= $(APP_NAME):latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

verify:: mod-verify errcheck

resolve-local:
	GO111MODULE=on go mod download

ensure-local:
	@echo "Go modules present in component - omitting."

dep-status:
	@echo "Go modules present in component - omitting."

dep-status-local:
	@echo "Go modules present in component - omitting."

mod-verify-local:
	GO111MODULE=on go mod verify 

# Run tests
test: export KUBEBUILDER_CONTROLPLANE_START_TIMEOUT = 2m
test: export KUBEBUILDER_CONTROLPLANE_STOP_TIMEOUT = 2m
test: generate verify manifests
	GO111MODULE=on go test ./... -coverprofile cover.out

# Install CRDs into a cluster
install: manifests
	kustomize build config/crd | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall: manifests
	kustomize build config/crd | kubectl delete -f -

removewebhook:
	-kubectl delete mutatingwebhookconfigurations.admissionregistration.k8s.io function-webhook.serverless.kyma-project.io

undeploy: uninstall removewebhook
	kustomize build config/default | ko delete  -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	cd config/manager && kustomize edit set image controller=${IMG}
	kustomize build config/default | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
.PHONY: manifests
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases


# Generate code
.PHONY: generate
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile=./hack/boilerplate.go.txt paths="./..."

# find or download controller-gen
# download controller-gen if necessary
.PHONY: controller-gen
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.4 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

.PHONY: path-to-referenced-charts
path-to-referenced-charts:
	@echo "resources/function-controller"

$(eval $(call buildpack-mount,resolve)) 
$(eval $(call buildpack-mount,mod-verify)) 
$(eval $(call buildpack-mount,test))
$(eval $(call buildpack-mount,controller-gen))
$(eval $(call buildpack-mount,generate))
$(eval $(call buildpack-mount,manifests))