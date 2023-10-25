package main

import (
	"flag"
	"time"

	"github.com/nice-pink/skupper-devops/pkg/autoheal"
)

func main() {
	// flags
	loop := flag.Bool("loop", false, "Loop check.")
	loopDelay := flag.Int("loopDelay", 10, "Loop check delay.")
	kubeConfig := flag.String("kubeConfig", ".kube/config", "ube config path.")
	restart := flag.Bool("restart", false, "Allow restarting servic if offline.")
	serviceRestarts := flag.Int("serviceRestarts", 10, "Max restart tries. Before giving up.")
	publishMetrics := flag.Bool("publishMetrics", false, "Publish metrics about auto-healing.")
	flag.Parse()

	// fmt.Println("--------")
	// fmt.Println("SRC: " + *src)
	// fmt.Println("DEST: " + *dest)
	// fmt.Println("--------")
	// fmt.Println("")

	autoheal.KubeConfigPath = *kubeConfig
	autoheal.MaxServiceRestarts = *serviceRestarts
	autoheal.PublishMetrics = *publishMetrics
	autoheal.Setup()

	// prepare
	for {
		autoheal.WatchServiceUptimes(*restart)

		if !*loop {
			break
		}

		time.Sleep(time.Duration(*loopDelay) * time.Second)
	}

}
