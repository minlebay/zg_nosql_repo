package shard_manager

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
	Tracer     trace.Tracer
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
		Tracer:     otel.Tracer("nosql-repo-shard-manager"),
	}
}

func (m *Manager) StartManager() {
	m.Logger.Info("Shard manager started")
}

func (m *Manager) StopManager() {
	m.wg.Wait()
	m.Logger.Info("Shard manager stopped")
}

func (m *Manager) Consume(ctx context.Context, msg *model.Message) {
	m.wg.Add(1)
	defer m.wg.Done()

	ctx, span := m.Tracer.Start(ctx, "Consume message")
	defer span.End()

	shardIndex, err := m.GetShardIndex(ctx, msg.Uuid)
	if err != nil {
		m.Logger.Error("Failed to get shard index", zap.Error(err))
		return
	}
	span.SetAttributes(attribute.Int("shard_index", shardIndex))

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
	dbsCount := len(m.Repository.DBs)
	if dbsCount == 0 {
		return 0, nil
	}
	uuidBytes := []byte(uuid)
	hash := crc32.ChecksumIEEE(uuidBytes)
	shardNumber := int(hash) % dbsCount
	return shardNumber, nil
}
