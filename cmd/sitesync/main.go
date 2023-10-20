package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/nice-pink/skupper-devops/pkg/sitesync"
)

func main() {
	configName := flag.String("configName", "", "Name of config map skupper-site to watch for changes.")
	configNamespace := flag.String("configNamespace", "", "Namespace of config map skupper-site to watch for changes.")
	isInCluster := flag.Bool("isInCluster", false, "Is running in cluster.")
	kubeConfig := flag.String("kubeConfig", "", "Optional. Default: .kube/config. If in cluster this field is unsed.")
	loopDelay := flag.Int("loopDelay", 10, "How many seconds between the runs?")
	loop := flag.Bool("loop", false, "Should keep on running?")
	initDelete := flag.Bool("initDelete", false, "Delete controllers on init.")
	flag.Parse()

	fmt.Println("ConfigMap")
	fmt.Println("name:", *configName)
	fmt.Println("namespace:", *configNamespace)
	fmt.Println("loop:", *loop)
	fmt.Println("initDelete:", *initDelete)

	if *configName == "" || *configNamespace == "" {
		flag.Usage()
		os.Exit(2)
	}

	sitesync.Setup(*configName, *configNamespace, *isInCluster, *kubeConfig)

	for {
		sitesync.Run(*initDelete)

		// should loop?
		if !*loop {
			break
		}

		time.Sleep(time.Duration(*loopDelay) * time.Second)
	}
}
