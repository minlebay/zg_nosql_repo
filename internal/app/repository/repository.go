package repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"sync"
	"time"
	"zg_nosql_repo/internal/model"
)

type Repository struct {
	Config     *Config
	Logger     *zap.Logger
	wg         sync.WaitGroup
	DB         *mongo.Database
	Collection *mongo.Collection
}

func NewRepository(logger *zap.Logger, config *Config) *Repository {
	return &Repository{
		Config: config,
		Logger: logger,
	}
}

func (r *Repository) StartRepository(ctx context.Context) {
	go func() {
		mongoURI := r.Config.Url
		databaseName := r.Config.DatabaseName

		dbctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		clientOptions := options.Client().ApplyURI(mongoURI)
		client, err := mongo.Connect(dbctx, clientOptions)
		if err != nil {
			r.Logger.Fatal("Failed to connect to MongoDB: %v", zap.Error(err))
		}
		defer func() {
			if err = client.Disconnect(dbctx); err != nil {
				r.Logger.Fatal("Failed to disconnect from MongoDB: %v", zap.Error(err))
			}
		}()
		r.DB = client.Database(databaseName)
	}()
}

func (r *Repository) StopRepository(ctx context.Context) {
	r.wg.Wait()
	r.Logger.Info("Repo started")
}

func (r *Repository) GetAll(ctx context.Context) ([]*model.Message, error) {
	r.Collection = r.DB.Collection("messages")

	var entities []*model.Message
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, model.NewAppErrorWithType(model.UnknownError)
	}
	if err := cursor.All(ctx, &entities); err != nil {
		return nil, model.NewAppErrorWithType(model.UnknownError)
	}
	return entities, nil
}

func (r *Repository) Create(ctx context.Context, entity *model.Message) (*model.Message, error) {
	r.Collection = r.DB.Collection("messages")

	_, err := r.Collection.InsertOne(ctx, entity)
	if err != nil {
		var mongoErr mongo.WriteException
		if errors.As(err, &mongoErr) {
			for _, writeErr := range mongoErr.WriteErrors {
				if writeErr.Code == 11000 {
					return nil, model.NewAppErrorWithType(model.ResourceAlreadyExists)
				}
			}
		}
		return nil, model.NewAppErrorWithType(model.UnknownError)
	}
	return entity, nil
}

func (r *Repository) GetById(ctx context.Context, id string) (*model.Message, error) {
	r.Collection = r.DB.Collection("messages")

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, model.NewAppErrorWithType(model.UnknownError)
	}
	var entity model.Message
	err = r.Collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&entity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, model.NewAppErrorWithType(model.NotFound)
		}
		return nil, model.NewAppErrorWithType(model.UnknownError)
	}
	return &entity, nil
}

func (r *Repository) Update(ctx context.Context, id string, entity *model.Message) (*model.Message, error) {
	r.Collection = r.DB.Collection("messages")

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, model.NewAppErrorWithType(model.UnknownError)
	}

	update := bson.M{"$set": entity}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := r.Collection.FindOneAndUpdate(ctx, bson.M{"_id": objId}, update, opts)
	if err = result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, model.NewAppErrorWithType(model.NotFound)
		}
		return nil, model.NewAppErrorWithType(model.UnknownError)
	}

	err = result.Decode(&entity)
	if err != nil {
		return nil, model.NewAppErrorWithType(model.UnknownError)
	}
	return entity, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	r.Collection = r.DB.Collection("messages")

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.NewAppErrorWithType(model.UnknownError)
	}
	res, err := r.Collection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return model.NewAppErrorWithType(model.UnknownError)
	}
	if res.DeletedCount == 0 {
		return model.NewAppErrorWithType(model.NotFound)
	}
	return nil
}
