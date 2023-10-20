package sitesync

import (
	"os"
	"time"

	"github.com/nice-pink/skupper-devops/pkg/kynetes"
	"github.com/nice-pink/skupper-devops/pkg/logger"
)

// Definitions

type SiteConfig struct {
	Name      string
	Namespace string
	Version   string
}

// Variabales

var (
	siteConfig       SiteConfig
	initiallyDeleted bool = false
)

// constants
const (
	SITE_CONTROLLER    string = "skupper-site-controller"
	SERVICE_CONTROLLER string = "skupper-service-controller"
	ROUTER             string = "skupper-router"
)

// Public

func Setup(name string, namespace string, isInCluster bool, kubeConfigPath string) {
	// kynetes
	kynetes.Verbose = false
	kynetes.IsInCluster = isInCluster
	if kubeConfigPath != "" {
		kynetes.KubeConfigPath = kubeConfigPath
	}

	if !isInCluster {
		// check if kube config exists
		if _, err := os.Stat(kynetes.KubeConfigPath); err != nil {
			logger.Log("Kube config file does not exist.")
			panic(err)
		}
	}

	// skupper config
	initConfig(name, namespace)
}

func Run(initDelete bool) {
	// get current version
	currentVersion := kynetes.ConfigMapGetResourceVersion(siteConfig.Name, siteConfig.Namespace)
	if currentVersion != "" {

		// if not the same restart controllers
		if (initDelete && !initiallyDeleted) || !compareCurrentVersion(currentVersion) {
			initiallyDeleted = true
			if restartControllers() {
				// update current version if successfully restarted controllers
				siteConfig.Version = currentVersion
			}
		}
	}
}

// init site

func initConfig(name string, namespace string) {
	siteConfig.Name = name
	siteConfig.Namespace = namespace
}

//

func compareCurrentVersion(currentVersion string) bool {
	if currentVersion == "" {
		logger.Error("Current version must not be empty!")
		return false
	}

	// set inital version
	if siteConfig.Version == "" {
		logger.Log("Initially set resource version:", currentVersion)
		siteConfig.Version = currentVersion
		return true
	}

	// compare versions
	if currentVersion != siteConfig.Version {
		logger.Log("Versions differ " + currentVersion + " != " + siteConfig.Version)
		return false
	}

	return true
}

func deleteController(name string, namespace string, retries int, delay int) bool {
	// delete controller
	for i := 0; i < retries; i++ {
		if err := kynetes.DeleteDeployment(name, namespace); err == nil {
			logger.Log("SUCCESS: Deleted", name, "in", namespace)
			return true
		}
		if i < 9 {
			time.Sleep(time.Duration(delay) * time.Second)
		}
		logger.Log("Retry deleting", name, "in", namespace)
	}
	return false
}

func controllerIsDeleted(name string, namespace string, retries int, delay int) bool {
	// wait until is deleted
	for i := 0; i < retries; i++ {
		if err := kynetes.GetDeployment(name, namespace); err != nil {
			logger.Log("SUCCESS:", name, "in", namespace, "removed entirely.")
			return true
		}
		if i < 9 {
			// exponential wait
			time.Sleep(time.Duration(2*i*delay) * time.Second)
		}
		logger.Log(name, "in", namespace, "still exists.")
	}
	return false
}

func restartControllers() bool {
	retries := 10
	delay := 1
	namespace := siteConfig.Namespace

	// delete service controller
	if !deleteController(SERVICE_CONTROLLER, namespace, retries, delay) {
		return false
	}
	// is service controller deleted?
	if !controllerIsDeleted(SERVICE_CONTROLLER, namespace, retries, delay) {
		return false
	}

	// router
	deletedRouter := deleteController(ROUTER, namespace, retries, delay)
	if !deletedRouter {
		// restart site controller anyways but at the end don't return success!
	} else {
		deletedRouter = controllerIsDeleted(ROUTER, namespace, retries, delay)
	}

	// needs delay here
	time.Sleep(time.Duration(5) * time.Second)

	// restart site controller
	for {
		restarted := kynetes.RestartDeployment(SITE_CONTROLLER, namespace)
		if restarted == nil {
			logger.Log("Restarted", SITE_CONTROLLER)
			break
		}
		logger.Log("Retry restarting", SITE_CONTROLLER)
	}

	// success ?
	return deletedRouter
}
