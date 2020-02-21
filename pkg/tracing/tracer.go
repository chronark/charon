package tracing

import (
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"io"
)

type LogrusAdapter struct{}

func (l LogrusAdapter) Error(msg string) {
	log.Errorf(msg)
}

func (l LogrusAdapter) Infof(msg string, args ...interface{}) {
	log.Infof(msg, args...)
}

func NewTracer(serviceName string) (opentracing.Tracer, io.Closer, error) {
	cfg := config.Configuration{
		ServiceName: serviceName,
		RPCMetrics:  true,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "jaeger:5775",
		},
	}
	tracer, closer, err := cfg.NewTracer(
		config.Logger(LogrusAdapter{}),
		config.Metrics(prometheus.New()),
	)
	return tracer, closer, err
}
