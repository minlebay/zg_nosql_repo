package nosql

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewModule() fx.Option {

	return fx.Module(
		"repo",
		fx.Provide(
			NewRepoConfig,
			NewRepo,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, r *Repo) {
				lc.Append(fx.StartStopHook(r.StartRepo, r.StopRepo))
			},
		),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("repo")
		}),
	)
}
