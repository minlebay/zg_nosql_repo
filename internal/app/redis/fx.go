package redis

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewModule() fx.Option {

	return fx.Module(
		"redis",
		fx.Provide(
			NewRedisConfig,
			NewRedis,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, r *Redis) {
				lc.Append(fx.StartStopHook(r.StartRedis, r.StartRedis))
			},
		),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("redis")
		}),
	)
}
