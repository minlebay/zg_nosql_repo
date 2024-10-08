package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
	"zg_nosql_repo/internal/app/shard_manager"
	"zg_nosql_repo/internal/model"
)

type Kafka struct {
	Config   *Config
	Logger   *zap.Logger
	Reader   *kafka.Reader
	SManager *shard_manager.Manager
	wg       sync.WaitGroup
}

func NewKafka(logger *zap.Logger, config *Config, sm *shard_manager.Manager) *Kafka {
	return &Kafka{
		Config:   config,
		Logger:   logger,
		SManager: sm,
	}
}

func (k *Kafka) StartKafka() {
	go func() {
		brokers := strings.Split(k.Config.Address, ",")
		k.Reader = kafka.NewReader(kafka.ReaderConfig{
			Brokers:          brokers,
			Topic:            k.Config.Topics,
			MinBytes:         10e3, // 10KB
			MaxBytes:         10e6, // 10MB
			ReadBatchTimeout: 1,
			GroupID:          k.Config.GroupID,
		})
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()
		for {
			k.Receive(ctx)
			time.Sleep(1 * time.Second)
		}
	}()
	k.Logger.Info("Kafka writer initialized", zap.String("address", k.Config.Address), zap.String("topic", k.Config.Topics))
}

func (k *Kafka) StopKafka() {
	k.wg.Wait()

	if err := k.Reader.Close(); err != nil {
		k.Logger.Error("Failed to close reader", zap.Error(err))
	} else {
		k.Logger.Info("Kafka reader closed successfully")
	}
}

func (k *Kafka) Receive(ctx context.Context) {
	k.wg.Add(1)
	defer k.wg.Done()

	msg := new(model.Message)
	m, err := k.Reader.ReadMessage(ctx)
	if err != nil {
		k.Logger.Error("Failed to read message", zap.Error(err))
	}

	err = msg.Unmarshal(m.Value)
	if err != nil {
		k.Logger.Error("Failed to unmarshal message", zap.Error(err))
		return
	}

	go k.SManager.Consume(ctx, msg)
}
