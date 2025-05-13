package metrics

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Collector struct {
	logger          *slog.Logger
	reg             *prometheus.Registry
	httpReqTotal    *prometheus.CounterVec
	httpReqDuration *prometheus.HistogramVec
}

func NewCollector(logger *slog.Logger) *Collector {
	col := &Collector{
		logger: logger.With(slog.String("name", "metrics")),
		reg:    prometheus.NewRegistry(),
		httpReqTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "blog",
			Name:      "http_requests_total",
			Help:      "Total number of HTTP requests.",
		},
			[]string{"method", "endpoint", "status"},
		),
		httpReqDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "blog",
			Name:      "request_duration_ms",
			Help:      "Histogram of request duration in milliseconds.",
			Buckets:   prometheus.DefBuckets,
		},
			[]string{"method", "endpoint"},
		),
	}
	col.reg.MustRegister(col.httpReqTotal)
	col.reg.MustRegister(col.httpReqDuration)

	return col
}

func (c *Collector) ServeMetrics(addr string) {
	mux := http.NewServeMux()
	mux.Handle("GET /metrics", promhttp.HandlerFor(c.reg, promhttp.HandlerOpts{}))

	c.logger.Info("starting metrics server", slog.String("addr", addr))
	if err := http.ListenAndServe(addr, mux); err != nil {
		c.logger.Error("failed to start metrics server", slog.String("err", err.Error()))
	}
}

func (c *Collector) IncHTTPRequest(method, url string, statusCode int) {
	c.httpReqTotal.WithLabelValues(method, url, strconv.Itoa(statusCode)).Inc()
}

func (c *Collector) RequestDuration(method, endpoint string, duration time.Duration) {
	c.httpReqDuration.WithLabelValues(method, endpoint).Observe(float64(duration.Milliseconds()))
}
