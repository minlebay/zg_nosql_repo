package nosql

import (
	"context"

	"go.uber.org/zap"
	"sync"
)

type Repo struct {
	Config *Config
	Logger *zap.Logger
	wg     sync.WaitGroup
}

func NewRepo() {

}

func (r *Repo) StartRepo(ctx context.Context) {
	go func() {}()
}

func (r *Repo) StopRepo(ctx context.Context) {
	r.wg.Wait()
	r.Logger.Info("Repo started")
}
