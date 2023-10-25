package autoheal

import (
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"github.com/nice-pink/skupper-devops/pkg/kynetes"
	"github.com/nice-pink/skupper-devops/pkg/logger"
	"github.com/nice-pink/skupper-devops/pkg/metric"
)

var (
	PrometheusUrl      string       = "http://localhost:9090"
	HttpClient         *http.Client = &http.Client{Timeout: 10 * time.Second}
	MaxServiceRestarts int          = 10
	Restarts                        = map[string]int{}
	PublishMetrics     bool         = false
)

type MetricSimple struct {
	Instance  string
	Namespace string
	Value     string
}

func Setup(kubeConfigPath string, isInCluster bool) {
	// setup kynetes
	kynetes.IsInCluster = isInCluster
	if kubeConfigPath != "" {
		kynetes.KubeConfigPath = kubeConfigPath
	}

	// setup metrics publisher
	if PublishMetrics {
		logger.Log("Publish metrics about auto-healing.")
		go func() {
			metric.Listen()
		}()
	}
}

func getJson(url string, target interface{}) error {
	resp, err := HttpClient.Get(url)
	if err != nil {
		logger.Error(err)
		panic(err)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func getSimpleMetrics(metrics []metric.Metric) []MetricSimple {
	var metricsSimple []MetricSimple

	pattern := `(http://)([a-zA-Z_-]*-skupper).([a-zA-Z_-]*).(.*)`
	for _, metric := range metrics {
		regex := regexp.MustCompile(pattern)
		value := regex.ReplaceAllString(metric.Info.Instance, "${2}")
		namespace := regex.ReplaceAllString(metric.Info.Instance, "${3}")
		metricSimple := &MetricSimple{value, namespace, metric.Value[1].(string)}
		metricsSimple = append(metricsSimple, *metricSimple)
	}

	return metricsSimple
}

func triggerAlert(instance string) {
	logger.Error("Alert for", instance)
}

func restartService(instance string, namespace string, restartCount int) bool {
	if restartCount == MaxServiceRestarts {
		triggerAlert(instance)
		return false
	} else if restartCount > MaxServiceRestarts {
		logger.Error(instance, "Restarts exeeded! Waiting for manual action.")
		return false
	}

	logger.Log("Restart", instance, "-", namespace)
	err := kynetes.DeleteService(instance, namespace)
	if err != nil {
		panic(err)
	}

	return true
}

func increaseRestart(instance string) {
	val, exists := Restarts[instance]
	if exists {
		Restarts[instance] = val + 1
	} else {
		Restarts[instance] = 1
	}
}

func WatchServiceUptimes(allowRestart bool) {
	^ := PrometheusUrl + metric.QueryPath + metric.UptimeMetricName
	response := new(metric.Response)
	getJson(url, response)

	//logger.Log(response.Data.Result[1].Info.Instance, response.Data.Result[1].Value[1])
	metricsSimple := getSimpleMetrics(response.Data.Result)
	for _, metric := range metricsSimple {
		if metric.Value == "0" {
			// found offline service
			logger.Error("Found service offline:", metric.Instance)
			increaseRestart(metric.Instance)

			if PublishMetrics {
				metricName := metric.Instance + "_" + metric.Namespace
				incCounter(metricName)
			}

			// restart
			if allowRestart {
				restartService(metric.Instance, metric.Namespace, Restarts[metric.Instance])
			}
		} else {
			// reset restarts if online
			Restarts[metric.Instance] = 0
		}
	}
}
