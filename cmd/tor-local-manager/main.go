package main

import (
	"flag"
	"log"

	configlib "github.com/kubernetes-sigs/kubebuilder/pkg/config"

	"github.com/kragniz/tor-controller/pkg/local"
)

// tor-manager main.
func main() {
	flag.Parse()

	//stopCh := signals.SetupSignalHandler()

	config := configlib.GetConfigOrDie()

	localManager := local.New(config)
	err := localManager.Run()
	if err != nil {
		log.Fatalf("%v", err)
	}
}
