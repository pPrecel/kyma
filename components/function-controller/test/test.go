package test

import (
	stdlog "log"

	"k8s.io/kubernetes/pkg/util/file"
)

func FileExists(arg string) {
	exists, err := file.FileExists(arg)
	if err != nil {
		stdlog.Fatal(err)
	}
	if !exists {
		stdlog.Fatalf("File %s does not exists", arg)
	}
}
