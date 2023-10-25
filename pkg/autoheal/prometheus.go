package autoheal

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	instanceCounters = map[string]prometheus.Counter{}
)

func incCounter(instance string) {
	// if counter exists -> inc
	val, ok := instanceCounters[instance]
	if ok {
		val.Inc()
		return
	}

	// init new counter for instance and inc
	counter := promauto.NewCounter(prometheus.CounterOpts{
		Name:        "skupper_service_auto_heal",
		Help:        "Skupper service was auto-healed.",
		ConstLabels: prometheus.Labels{"instance": instance},
	})
	counter.Inc()

	// save to map
	instanceCounters[instance] = counter
}
