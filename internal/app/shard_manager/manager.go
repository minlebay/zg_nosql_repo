package shard_manager

import (
	"context"
	"go.uber.org/zap"
	"sync"
	"zg_nosql_repo/internal/app/redis"
	"zg_nosql_repo/internal/app/repository"
	"zg_nosql_repo/internal/model"
)

type Manager struct {
	Config     *Config
	Logger     *zap.Logger
	Redis      *redis.Redis
	Repository *repository.Repository
	Messages   chan *model.Message
	wg         sync.WaitGroup
}

func NewManager(
	logger *zap.Logger,
	config *Config,
	redis *redis.Redis,
	repo *repository.Repository,
) *Manager {
	return &Manager{
		Config:     config,
		Logger:     logger,
		Redis:      redis,
		Repository: repo,
	}
}

func (m *Manager) StartManager(ctx context.Context) {
	//go func() {
	//	for {
	//		select {
	//		case msg := <-m.Messages:
	//			m.Logger.Info("Message received", zap.String("uuid", msg.UUID))
	//			m.consume(context.Background(), msg)
	//		default:
	//			continue
	//		}
	//	}
	//}()
	m.Logger.Info("Shard manager started")
}

func (m *Manager) StopManager(ctx context.Context) {
	m.wg.Wait()
	m.Logger.Info("Shard manager stopped")
}

// Calculate shard number and send message to the shard
// Store index to the redis
func (m *Manager) Consume(ctx context.Context, msg *model.Message) {
	m.wg.Add(1)
	defer m.wg.Done()

	m.Logger.Info("Message stored", zap.String("uuid", msg.Uuid))

	//shard := m.Repository.GetShard(msg.Uuid)
	//m.Logger.Info("Message received", zap.String("uuid", msg.Uuid), zap.Int("shard", shard))
	//
	//err := m.Repository.StoreIndex(ctx, msg.Uuid, shard)
	//if err != nil {
	//	m.Logger.Error("Failed to store index", zap.Error(err))
	//}
}
