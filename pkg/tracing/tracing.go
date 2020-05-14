package tracing

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"runtime"
)

func StartSpanFromContext(ctx context.Context) (opentracing.Span, context.Context) {
	span, traceContext := opentracing.StartSpanFromContext(ctx, getCallerName(0))
	return span, traceContext
}

func StartSpanFromContextWithTracer(ctx context.Context, tracer opentracing.Tracer, name string) (opentracing.Span, context.Context) {
	return opentracing.StartSpanFromContextWithTracer(ctx, tracer, name)
}

func StartSpanFromContextf(ctx context.Context, name string, args ...interface{}) (opentracing.Span, context.Context) {
	span, traceContext := opentracing.StartSpanFromContext(ctx, getCallerName(0)+" "+fmt.Sprintf(name, args...))
	return span, traceContext
}

func StartNamedSpanFromContext(ctx context.Context, name string) (opentracing.Span, context.Context) {
	span, traceContext := opentracing.StartSpanFromContext(ctx, name)
	return span, traceContext
}

func getCallerName(offset int) string {
	pc, _, _, ok := runtime.Caller(2 + offset)
	if !ok {
		return "<unknown>"
	}
	details := runtime.FuncForPC(pc)
	return fmt.Sprintf("%s()", details.Name())
}
