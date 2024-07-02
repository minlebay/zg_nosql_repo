package shard_manager

import (
	"context"
	"go.uber.org/zap"
	"hash/crc32"
	"strconv"
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
	m.Logger.Info("Shard manager started")
}

func (m *Manager) StopManager(ctx context.Context) {
	m.wg.Wait()
	m.Logger.Info("Shard manager stopped")
}

func (m *Manager) Consume(ctx context.Context, msg *model.Message) {
	m.wg.Add(1)
	defer m.wg.Done()

	// Calculate shard number
	shardIndex, err := m.GetShardIndex(ctx, msg.Uuid)
	if err != nil {
		m.Logger.Error("Failed to get shard index", zap.Error(err))
		return
	}

	created, err := m.Repository.Create(ctx, *m.Repository.DBs[shardIndex], msg)
	if err != nil {
		m.Logger.Error("Failed to store message", zap.Error(err))
		return
	}
	m.Logger.Info("Message stored", zap.String("uuid", created.Uuid), zap.Int("shard", shardIndex))

	bytes := []byte(strconv.Itoa(shardIndex))
	err = m.Redis.Put(created.Uuid, bytes)
	if err != nil {
		m.Logger.Error("Failed to store index", zap.Error(err))
		return
	}
	m.Logger.Info("Index stored", zap.String("uuid", created.Uuid), zap.Int("shard", shardIndex))
}

func (m *Manager) GetShardIndex(ctx context.Context, uuid string) (int, error) {
	uuidBytes := []byte(uuid)
	hash := crc32.ChecksumIEEE(uuidBytes)
	shardNumber := int(hash) % len(m.Repository.DBs)
	return shardNumber, nil
}
