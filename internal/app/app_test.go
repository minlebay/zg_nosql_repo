package app

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"testing"
	"zg_nosql_repo/internal/app/kafka"
	"zg_nosql_repo/internal/app/log"
	"zg_nosql_repo/internal/app/redis"
	"zg_nosql_repo/internal/app/repository"
	"zg_nosql_repo/internal/app/shard_manager"
	"zg_nosql_repo/internal/app/tracer"
)

func TestValidateApp(t *testing.T) {
	err := fx.ValidateApp(
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
	require.NoError(t, err)
}
