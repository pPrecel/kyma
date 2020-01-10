package webhook

import (
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// Add adds itself to the manager
func Add(mgr manager.Manager) error {
	log = logf.Log.WithName("webhook_server")
	log.Info("setting up webhook server")
	srv := mgr.GetWebhookServer()
	srv.CertDir = "/tmp/cert"
	srv.Port = 9876
	srv.Register("/"+"mutating-create-function", &webhook.Admission{Handler: &FunctionCreateHandler{}})
	return nil
}
