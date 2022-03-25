package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fitant/xbin-api/src/db/internal/migrations"
	"github.com/fitant/xbin-api/src/utils"

	"github.com/fitant/xbin-api/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoFindInsert struct {
	ephSnips *mongo.Collection
	snips    *mongo.Collection
	timeout  time.Duration
}

var ErrNoDocuments error = mongo.ErrNoDocuments
var ErrDuplicateKey error = errors.New("error: duplicate key insertion")

func NewMongoStore(cfg *config.Config) (*MongoFindInsert, error) {
	timeout := time.Duration(cfg.DB.TimeoutInSec() * int(time.Second))
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DB.Address()))
	if err != nil {
		panic(fmt.Sprintf("%s : %v", "[DB] [NewMongoStore] [Connect]", err))
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		panic(fmt.Sprintf("%s : %v", "[DB] [NewMongoStore] [Ping]", err))
	}

	db := client.Database(cfg.DB.Database())

	migrations.Migrate(db.Client(), cfg)

	return &MongoFindInsert{
		snips:    db.Collection(cfg.DB.Collection()),
		ephSnips: db.Collection(cfg.DB.EphemeralCollection().Name),
		timeout:  timeout,
	}, nil
}

func (db *MongoFindInsert) FindOne(condition []byte, eph bool) (*mongo.SingleResult, error) {
	collection := db.snips
	if eph {
		collection = db.ephSnips
	}

	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	res := collection.FindOne(ctx, condition)
	if res.Err() != nil && res.Err() != mongo.ErrNoDocuments {
		utils.Logger.Error(fmt.Sprintf("%s : %s", "[DB] [FindOne]", res.Err()))
		return nil, res.Err()
	}

	if res.Err() == mongo.ErrNoDocuments {
		return nil, ErrNoDocuments
	}

	return res, nil
}

func (db *MongoFindInsert) InsertOne(document []byte, eph bool) (*mongo.InsertOneResult, error) {
	collection := db.snips
	if eph {
		collection = db.ephSnips
	}

	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	res, err := collection.InsertOne(ctx, document)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, ErrDuplicateKey
		}
		utils.Logger.Error(fmt.Sprintf("%s : %v", "[DB] [InsertOne]", err))
		return nil, err
	}

	return res, nil
}
