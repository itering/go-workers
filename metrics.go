package workers

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	TaskProcessed = prometheus.NewCounter(prometheus.CounterOpts{
		Subsystem: "worker_task",
		Name:      "processed",
		Help:      "Total number of tasks processed.",
	})

	TaskFailed = prometheus.NewCounter(prometheus.CounterOpts{
		Subsystem: "worker_task",
		Name:      "failed",
		Help:      "Total number of tasks failed.",
	})

	TaskRetries = prometheus.NewCounter(prometheus.CounterOpts{
		Subsystem: "worker_task",
		Name:      "retries",
		Help:      "Total number of tasks retried.",
	})

	TasksInqueue = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "worker_task_inqueue",
			Help: "number of tasks in queue.",
		},
		[]string{"type"},
	)
)

func init() {
	prometheus.MustRegister(TaskProcessed, TaskFailed, TaskRetries, TasksInqueue)
}

func MetricsServer(port int) {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
