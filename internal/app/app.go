package app

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"zg_nosql_repo/internal/app/kafka"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Options(
			kafka.NewModule(),
		),
		fx.Provide(
			zap.NewProduction,
			NewConfig,
		),
	)
}
