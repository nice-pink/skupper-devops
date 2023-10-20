package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nice-pink/skupper-devops/pkg/sitesync"
)

func main() {
	fmt.Println("main")

	configName := flag.String("configName", "", "Name of config map skuppter-site to watch for changes.")
	configNamespace := flag.String("configNamespace", "", "Namespace of config map skuppter-site to watch for changes.")
	isInCluster := flag.Bool("isInCluster", false, "Is running in cluster.")
	kubeConfig := flag.String("kubeConfig", "", "Optional. Default: .kube/config. If in cluster this field is unsed.")
	flag.Parse()

	fmt.Println("ConfigMap")
	fmt.Println("name:", *configName)
	fmt.Println("namespace:", *configNamespace)

	if *configName == "" || *configNamespace == "" {
		flag.Usage()
		os.Exit(2)
	}

	sitesync.Setup(*configName, *configNamespace, *isInCluster, *kubeConfig)
	sitesync.Run()
}
