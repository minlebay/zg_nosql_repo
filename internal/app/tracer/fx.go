package tracer

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewModule() fx.Option {

	return fx.Module(
		"sharding",
		fx.Provide(
			NewTracerConfig,
			NewTracer,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, t *Tracer) {
				lc.Append(fx.StartStopHook(t.StartTracer, t.StopTracer))
			},
		),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("nosql_sharding")
		}),
	)
}
