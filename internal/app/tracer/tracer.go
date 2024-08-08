package tracer

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/zap"
)

type Tracer struct {
	Config *Config
	Logger *zap.Logger
	tp     *trace.TracerProvider
}

func NewTracer(config *Config, logger *zap.Logger) *Tracer {
	return &Tracer{
		Config: config,
		Logger: logger,
	}
}

func (t *Tracer) StartTracer() {
	go func() {
		url := t.Config.Url

		exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
		if err != nil {
			t.Logger.Error("Failed to create Jaeger exporter", zap.Error(err))
		}

		t.tp = trace.NewTracerProvider(
			trace.WithBatcher(exp),
			trace.WithResource(resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("zg_nosql_repo"),
			)),
		)
		otel.SetTracerProvider(t.tp)
	}()
}

func (t *Tracer) StopTracer() {
	t.tp.Shutdown(context.Background())
}
