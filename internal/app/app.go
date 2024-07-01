package app

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"zg_nosql_repo/internal/app/kafka"
	"zg_nosql_repo/internal/app/redis"
	"zg_nosql_repo/internal/app/repository"
	"zg_nosql_repo/internal/app/shard_manager"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Options(
			kafka.NewModule(),
			redis.NewModule(),
			repository.NewModule(),
			shard_manager.NewModule(),
		),
		fx.Provide(
			zap.NewProduction,
			NewConfig,
		),
	)
}
