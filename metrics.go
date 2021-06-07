package workers

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	TaskProcessed = prometheus.NewCounter(prometheus.CounterOpts{
		Subsystem: "",
		Name:      "worker_task_processed",
		Help:      "Total number of tasks processed.",
	})

	TaskFailed = prometheus.NewCounter(prometheus.CounterOpts{
		Subsystem: "",
		Name:      "worker_task_failed",
		Help:      "Total number of tasks failed.",
	})

	TaskRetries = prometheus.NewCounter(prometheus.CounterOpts{
		Subsystem: "",
		Name:      "worker_task_retries",
		Help:      "Total number of tasks retried.",
	})

	TaskEnqueue = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "worker_task_enqueue",
			Help: "number of task enqueue, grouped by queue name",
		},
		[]string{"name"},
	)

	TaskDequeue = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "worker_task_dequeue",
			Help: "number of task dequeue, grouped by queue name",
		},
		[]string{"name"},
	)
)

func init() {
	prometheus.MustRegister(TaskProcessed, TaskFailed, TaskRetries, TaskEnqueue, TaskDequeue)
}

func MetricsServer(port int) {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
