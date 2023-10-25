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
	kubeConfig := flag.String("kubeConfig", ".kube/config", "Kube config path.")
	isInCluster := flag.Bool("isInCluster", false, "Is running in cluster.")
	restart := flag.Bool("restart", false, "Allow restarting servic if offline.")
	serviceRestarts := flag.Int("serviceRestarts", 10, "Max restart tries. Before giving up.")
	publishMetrics := flag.Bool("publishMetrics", false, "Publish metrics about auto-healing.")
	prometheusUrl := flag.String("prometheusUrl", "", "Prometheus url.")
	flag.Parse()

	// fmt.Println("--------")
	// fmt.Println("SRC: " + *src)
	// fmt.Println("DEST: " + *dest)
	// fmt.Println("--------")
	// fmt.Println("")

	autoheal.MaxServiceRestarts = *serviceRestarts
	autoheal.PublishMetrics = *publishMetrics
	autoheal.PrometheusUrl = *prometheusUrl
	autoheal.Setup(*kubeConfig, *isInCluster)

	// prepare
	for {
		autoheal.WatchServiceUptimes(*restart)

		if !*loop {
			break
		}

		time.Sleep(time.Duration(*loopDelay) * time.Second)
	}

}
