package main

import (
	"flag"
	"time"

	"github.com/nice-pink/skupper-devops/pkg/uptime"
)

func main() {
	// flags
	loop := flag.Bool("loop", false, "Loop check.")
	loopDelay := flag.Int("loopDelay", 10, "Loop check delay.")
	kubeConfig := flag.String("kubeConfig", ".kube/config", "ube config path.")
	flag.Parse()

	// fmt.Println("--------")
	// fmt.Println("SRC: " + *src)
	// fmt.Println("DEST: " + *dest)
	// fmt.Println("--------")
	// fmt.Println("")

	uptime.KubeConfigPath = *kubeConfig
	uptime.Setup()

	// prepare
	for {
		uptime.WatchServiceUptimes()

		if !*loop {
			break
		}

		time.Sleep(time.Duration(*loopDelay) * time.Second)
	}

}
