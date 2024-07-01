package shard_manager

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewModule() fx.Option {

	return fx.Module(
		"sharding",
		fx.Provide(
			NewManagerConfig,
			NewManager,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, r *Manager) {
				lc.Append(fx.StartStopHook(r.StartManager, r.StartManager))
			},
		),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("sharding")
		}),
	)
}
