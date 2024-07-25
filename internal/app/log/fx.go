package log

import (
	"go.uber.org/fx"
)

func NewModule() fx.Option {

	return fx.Module(
		"logstash",
		fx.Provide(
			NewLogstashConfig,
			NewLogstashWriter,
			NewZapLogger,
		),
	)
}
