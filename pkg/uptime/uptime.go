package uptime

import (
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"github.com/nice-pink/skupper-devops/pkg/logger"
)

const (
	SERVICE_RESTARTS int = 4
)

var (
	PrometheusUrl string       = "http://localhost:9090"
	QueryPath     string       = "/api/v1/query?query="
	UptimeQuery   string       = "probe_success"
	HttpClient    *http.Client = &http.Client{Timeout: 10 * time.Second}
	Restarts                   = map[string]int{}
)

type MetricSimple struct {
	Instance string
	Value    string
}

func getJson(url string, target interface{}) error {
	resp, err := HttpClient.Get(url)
	if err != nil {
		logger.Error(err)
		panic(err)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func getSimpleMetrics(metrics []Metric) []MetricSimple {
	var metricsSimple []MetricSimple

	pattern := `(http://)([a-zA-Z_-]*)(.*)`
	for _, metric := range metrics {
		regex := regexp.MustCompile(pattern)
		value := regex.ReplaceAllString(metric.Info.Instance, "${2}")
		metricSimple := &MetricSimple{value, metric.Value[1].(string)}
		metricsSimple = append(metricsSimple, *metricSimple)
	}

	return metricsSimple
}

func triggerAlert(instance string) {
	logger.Error("Alert for", instance)
}

func restartService(instance string, restartCount int) bool {
	if restartCount == SERVICE_RESTARTS {
		triggerAlert(instance)
		return false
	} else if restartCount > SERVICE_RESTARTS {
		logger.Error(instance, "Restarts exeeded! Waiting for manual action.")
		return false
	}

	logger.Log("Restart", instance)
	return true
}

func increaseRestart(instance string) {
	val, exists := Restarts[instance]
	if exists {
		Restarts[instance] = val + 1
	} else {
		Restarts[instance] = 1
	}
	restartService(instance, Restarts[instance])
}

func WatchServiceUptimes() {
	url := PrometheusUrl + QueryPath + UptimeQuery
	response := new(Response)
	getJson(url, response)

	//logger.Log(response.Data.Result[1].Info.Instance, response.Data.Result[1].Value[1])
	metricsSimple := getSimpleMetrics(response.Data.Result)
	for _, metric := range metricsSimple {
		if metric.Value == "0" {
			logger.Error("Found service offline:", metric.Instance)
			increaseRestart(metric.Instance)
		} else {
			Restarts[metric.Instance] = 0
		}
	}
}

// func WatchServiceUptimes() {
// 	resp, err := http.Get()
// 	if err != nil {
// 		logger.Error(err)
// 		panic(err)
// 	}

// 	body, err := io.ReadAll(resp.Body)
// 	stringBody := string(body[:])
// 	logger.Log(stringBody)
// }
