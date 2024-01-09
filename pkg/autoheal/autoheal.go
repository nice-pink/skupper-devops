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
	ServiceReset       int          = 30
	Restarts                        = map[string]int{}
	PublishMetrics     bool         = false
)

type MetricSimple struct {
	Instance  string
	Namespace string
	Value     string
}

func Setup(kubeConfigPath string, isInCluster bool, metricsScrapePort int) {
	// setup kynetes
	kynetes.IsInCluster = isInCluster
	if kubeConfigPath != "" {
		kynetes.KubeConfigPath = kubeConfigPath
	}

	// setup metrics publisher
	if PublishMetrics {
		logger.Log("Publish metrics about auto-healing.")
		go func() {
			metric.Listen(metricsScrapePort)
		}()
	}
}

func getJson(url string, target interface{}) error {
	resp, err := HttpClient.Get(url)
	if err != nil {
		logger.Error(err)
		return err
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

	servicesGauge.Set(float64(len(metricsSimple)))

	return metricsSimple
}

func triggerAlert(instance string, namespace string) {
	logger.Error("Alert for", instance)
	incAlertsCounter(instance, namespace)
}

func restartService(instance string, namespace string, restartCount int) (restarted bool, reset bool) {
	if restartCount == MaxServiceRestarts {
		triggerAlert(instance, namespace)
		return false, false
	} else if restartCount > MaxServiceRestarts {
		if restartCount > ServiceReset {
			logger.Error(instance, "Reset restart counter.")
			return false, true
		}
		logger.Error(instance, "Restarts exeeded! Waiting for manual action.")
		return false, false
	}

	logger.Log("Restart", instance, "-", namespace)
	err := kynetes.DeleteService(instance, namespace)
	if err != nil {
		logger.Log(err)
		incServiceMissingCounter(instance, namespace)
		return false, false
	}

	return true, false
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
	// increment general counter
	checksCounter.Inc()

	// get prometheus metrics
	url := PrometheusUrl + metric.QueryPath + metric.UptimeMetricName
	response := new(metric.Response)
	err := getJson(url, response)
	if err != nil {
		return
	}

	//logger.Log(response.Data.Result[1].Info.Instance, response.Data.Result[1].Value[1])
	metricsSimple := getSimpleMetrics(response.Data.Result)

	for _, metric := range metricsSimple {
		if metric.Value == "0" {
			// found offline service
			logger.Error("Found service offline:", metric.Instance)
			increaseRestart(metric.Instance)

			incAutohealCounter(metric.Instance, metric.Namespace)

			// restart
			if allowRestart {
				_, reset := restartService(metric.Instance, metric.Namespace, Restarts[metric.Instance])
				if reset {
					Restarts[metric.Instance] = 0
				}
			}
		} else {
			// reset restarts if online
			Restarts[metric.Instance] = 0
		}
	}
}
