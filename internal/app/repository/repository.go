package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	url2 "net/url"
	"strings"
	"sync"
	"time"
	"zg_nosql_repo/internal/model"
)

type Repository struct {
	Config           *Config
	Logger           *zap.Logger
	wg               sync.WaitGroup
	DBs              []*mongo.Database
	Collection       *mongo.Collection
	CancelFunc       context.CancelFunc
	ClientDisconnect func()
}

func NewRepository(logger *zap.Logger, config *Config) *Repository {
	return &Repository{
		Config: config,
		Logger: logger,
	}
}

func (r *Repository) StartRepository() {
	go func() {
		for _, db := range r.Config.Dbs {
			url, err := url2.Parse(db)
			if err != nil {
				r.Logger.Fatal("Failed to parse MongoDB URL: %v", zap.Error(err))
			}
			databaseName := strings.TrimPrefix(url.Path, "/")

			dbctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			r.CancelFunc = cancel

			clientOptions := options.Client().ApplyURI(db)
			client, err := mongo.Connect(dbctx, clientOptions)
			if err != nil {
				r.Logger.Fatal("Failed to connect to MongoDB: %v", zap.Error(err))
			}
			r.ClientDisconnect = func() {
				if err = client.Disconnect(dbctx); err != nil {
					r.Logger.Fatal("Failed to disconnect from MongoDB: %v", zap.Error(err))
				}
			}
			r.DBs = append(r.DBs, client.Database(databaseName))
		}
	}()
}

func (r *Repository) StopRepository() {
	r.wg.Wait()
	r.ClientDisconnect()
	r.CancelFunc()

	r.Logger.Info("Repo started")
}

func (r *Repository) Create(ctx context.Context, db mongo.Database, entity *model.Message) (*model.Message, error) {
	r.Collection = db.Collection("messages")

	_, err := r.Collection.InsertOne(ctx, entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}
