package middlewares

import (
	"github.com/labstack/echo-contrib/prometheus"
)

func NewEchoPrometheusMiddleware() *prometheus.Prometheus {
	return prometheus.NewPrometheus("echo", nil)
}
