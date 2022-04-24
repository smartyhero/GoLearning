package metrics

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	MetricsNamespace = "smartyhero"
)

var (
	functionLatency = CreateExecutionTimeMetric(MetricsNamespace, "Time spent.")
)

func NewTimer() *ExecutionTimer {
	return NewExecutionTimer(functionLatency)
}
func Register() {
	err := prometheus.Register(functionLatency)
	if err != nil {
		fmt.Println(err)
	}
}

func NewExecutionTimer(histo *prometheus.HistogramVec) *ExecutionTimer {
	now := time.Now()
	return &ExecutionTimer{
		histo: histo,
		start: now,
		last:  now,
	}
}

func (t *ExecutionTimer) ObserveTotal() {
	(*t.histo).WithLabelValues("total").Observe(time.Since(t.start).Seconds())
}

func CreateExecutionTimeMetric(namespace string, help string) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "execution_latency_seconds",
			Help:      help,
			Buckets:   prometheus.ExponentialBuckets(0.001, 2, 15),
		}, []string{"step"},
	)
}

type ExecutionTimer struct {
	histo *prometheus.HistogramVec
	start time.Time
	last  time.Time
}
