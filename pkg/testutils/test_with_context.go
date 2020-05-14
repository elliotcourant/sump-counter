package testutils

import (
	"context"
	"flag"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/require"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"testing"
)

var (
	traceTests = flag.Bool("trace-tests", true, "enables jaeger tracing for unit tests")
)

func NewTestWithContext(t *testing.T, test func(t *testing.T, ctx context.Context)) {
	shouldTrace := traceTests != nil && *traceTests
	if shouldTrace {
		// Sample configuration for testing. Use constant sampling to sample every trace
		// and enable LogSpan to log every span via configured Logger.
		cfg := jaegercfg.Configuration{
			ServiceName: "Sump Boi",
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &jaegercfg.ReporterConfig{
				LogSpans: true,
			},
		}

		// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
		// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
		// frameworks.
		jLogger := jaegerlog.NullLogger
		jMetricsFactory := metrics.NullFactory

		// Initialize tracer with a logger and a metrics factory
		tracer, closer, err := cfg.NewTracer(
			jaegercfg.Logger(jLogger),
			jaegercfg.Metrics(jMetricsFactory),
		)
		require.NoError(t, err)

		// Set the singleton opentracing.Tracer with the Jaeger tracer.
		opentracing.SetGlobalTracer(tracer)
		defer closer.Close()
	}

	span, ctx := opentracing.StartSpanFromContext(context.Background(), t.Name())
	defer span.Finish()

	test(t, ctx)
}
