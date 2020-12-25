// +build prometheus

package instrumentation

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type promCounter struct {
	*prometheus.CounterVec
}

func NewCounter(name, help string, labels ...string) Counter {
	return &promCounter{
		promauto.NewCounterVec(prometheus.CounterOpts{
			Name: name,
			Help: help,
		}, labels),
	}
}

func (s *promCounter) Inc(values ...interface{}) {
	s.WithLabelValues(mapIntfToStringArray(values)...).Inc()
}

func mapIntfToStringArray(items []interface{}) []string {
	sArr := make([]string, len(items))
	for i, v := range items {
		sArr[i] = fmt.Sprintf("%v", v)
	}
	return sArr
}
