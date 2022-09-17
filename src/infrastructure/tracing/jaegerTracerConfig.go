package tracing

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

type JaegerTracerConfig struct {
	Tracer       *opentracing.Tracer
	TracerCloser io.Closer
}

func NewJaegerTracerConfig() *JaegerTracerConfig {
	host := os.Getenv("JAEGER_AGENT_HOST")
	port := os.Getenv("JAEGER_AGENT_PORT")
	cfg := config.Configuration{
		ServiceName: "uaa-tracer",
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  fmt.Sprintf("%s:%s", host, port),
		},
	}
	tracer, tracerCloser, err := cfg.NewTracer()
	if err != nil {
		panic("Could not initialize jaeger tracer: " + err.Error())
	}

	opentracing.SetGlobalTracer(tracer)
	return &JaegerTracerConfig{
		Tracer:       &tracer,
		TracerCloser: tracerCloser,
	}
}
