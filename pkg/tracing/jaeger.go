package tracing

import (
	"fmt"
	"github.com/chronark/charon/pkg/log"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/rpcmetrics"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"go.uber.org/zap"
	"io"
)

type jaegerLoggerAdapter struct {
	logger log.Logger
}

func (l jaegerLoggerAdapter) Error(msg string) {
	l.logger.Error(msg)
}

func (l jaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, args...))
}

func NewTracer(serviceName string, logger log.Factory) (opentracing.Tracer, io.Closer) {
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
	jaegerLogger := jaegerLoggerAdapter{logger.Bg()}
	metrics := prometheus.New()

	tracer, closer, err := cfg.NewTracer(
		config.Logger(jaegerLogger),
		config.Metrics(metrics),
		config.Observer(rpcmetrics.NewObserver(metrics, rpcmetrics.DefaultNameNormalizer)),
	)
	if err != nil {
		logger.Bg().Fatal("cannot initialize Jaeger Tracer", zap.Error(err))
	}
	return tracer, closer
}
