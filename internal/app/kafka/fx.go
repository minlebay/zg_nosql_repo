package kafka

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewModule() fx.Option {

	return fx.Module(
		"kafka",
		fx.Provide(
			NewKafkaConfig,
			NewKafka,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, r *Kafka) {
				lc.Append(fx.StartStopHook(r.StartKafka, r.StopKafka))
			},
		),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("kafka")
		}),
	)
}
