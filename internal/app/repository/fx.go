package repository

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewModule() fx.Option {

	return fx.Module(
		"repo",
		fx.Provide(
			NewRepositoryConfig,
			NewRepository,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, r *Repository) {
				lc.Append(fx.StartStopHook(r.StartRepository, r.StartRepository))
			},
		),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("repo")
		}),
	)
}
