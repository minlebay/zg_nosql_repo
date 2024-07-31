package app

import (
	"go.uber.org/fx"
	"zg_nosql_repo/internal/app/kafka"
	"zg_nosql_repo/internal/app/log"
	"zg_nosql_repo/internal/app/redis"
	"zg_nosql_repo/internal/app/repository"
	"zg_nosql_repo/internal/app/shard_manager"
	"zg_nosql_repo/internal/app/tracer"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Options(
			kafka.NewModule(),
			redis.NewModule(),
			repository.NewModule(),
			shard_manager.NewModule(),
			log.NewModule(),
			tracer.NewModule(),
		),
		fx.Provide(
			NewConfig,
		),
	)
}
