package autoheal

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	autohealCounters       = map[string]prometheus.Counter{}
	serviceMissingCounters = map[string]prometheus.Counter{}
	alertsCounters         = map[string]prometheus.Counter{}
	checksCounter          = promauto.NewCounter(prometheus.CounterOpts{
		Name: "skupper_service_checks",
		Help: "Count checks.",
	})
	servicesGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "skupper_services",
		Help: "Amount of services known to skupper.",
	})
)

func GetIdentifier(instance string, namespace string) string {
	return instance + "_" + namespace
}

func incAutohealCounter(instance string, namespace string) {
	identifier := GetIdentifier(instance, namespace)

	// if counter exists -> inc
	val, ok := autohealCounters[identifier]
	if ok {
		val.Inc()
		return
	}

	// init new counter for instance and inc
	counter := promauto.NewCounter(prometheus.CounterOpts{
		Name:        "skupper_service_auto_heal",
		Help:        "Skupper service was auto-healed.",
		ConstLabels: prometheus.Labels{"instance": identifier},
	})
	counter.Inc()

	// save to map
	autohealCounters[identifier] = counter
}

func incServiceMissingCounter(instance string, namespace string) {
	identifier := GetIdentifier(instance, namespace)

	// if counter exists -> inc
	val, ok := serviceMissingCounters[identifier]
	if ok {
		val.Inc()
		return
	}

	// init new counter for instance and inc
	counter := promauto.NewCounter(prometheus.CounterOpts{
		Name:        "skupper_service_missing",
		Help:        "Skupper service is missing.",
		ConstLabels: prometheus.Labels{"instance": identifier},
	})
	counter.Inc()

	// save to map
	serviceMissingCounters[identifier] = counter
}

func incAlertsCounter(instance string, namespace string) {
	identifier := GetIdentifier(instance, namespace)

	// if counter exists -> inc
	val, ok := alertsCounters[identifier]
	if ok {
		val.Inc()
		return
	}

	// init new counter for instance and inc
	counter := promauto.NewCounter(prometheus.CounterOpts{
		Name:        "skupper_service_alerts",
		Help:        "Alert triggered.",
		ConstLabels: prometheus.Labels{"instance": identifier},
	})
	counter.Inc()

	// save to map
	alertsCounters[identifier] = counter
}
