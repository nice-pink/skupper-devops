package main

import (
	"flag"
	"fmt"

	"github.com/nice-pink/skupper-devops/pkg/kynetes"
)

func main() {
	// flags
	kubeConfig := flag.String("kubeConfig", "", "Path to .kube/config")
	flag.Parse()

	fmt.Println("--------")
	fmt.Println("KUBE_CONFIG: " + *kubeConfig)
	fmt.Println("--------")
	fmt.Println("")

	// prepare
	_, err := kynetes.ReadResource("secret.yaml")
	if err != nil {
		panic(err)
	}

	// set kube config path
	kynetes.KubeConfigPath = *kubeConfig

	// get pods
	// kynetes.ListPods("streaming")
	// kynetes.GetPod("processing-engine-podcasts-78cdfd6fdb-2wkpf", "streaming")

	// create secret
	// secretByte, err := kynetes.ReadResource("secret.yaml")
	// if err != nil {
	// 	panic(err)
	// }
	// err = kynetes.CreateSecret("ops", secretByte)
	// if err != nil {
	// 	panic(err)
	// }

	// get secrets
	// kynetes.ListSecrets("ops")
	// sec := kynetes.GetSecret("test-secret", "ops", false)
	// fmt.Println(sec)

	// has data
	if kynetes.SecretHasData("test-secret", "ops") {
		fmt.Println("Has data")
	} else {
		fmt.Println("is empty")
	}

	// delete secret
	// err = kynetes.DeleteSecret("test-secret", "ops")
	// if err != nil {
	// 	panic(err)
	// }
}
