package mongodb

import (
	"btc-project/library/log"
	xtime "btc-project/library/time"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Config struct {
	Addr    		string 				`toml:"addr"`
	Timeout 		xtime.Duration    	`toml:"timeout"`
	DBName			string				`toml:"databaseName"`
}

type DB struct {
	Ctx context.Context
	Client *mongo.Client
	Database *mongo.Database
}

func NewMongoDb(c *Config) (db *DB) {
	if c.Addr == "" {
		panic("mongodb address not defined")
	}

	if c.Timeout <= 0 {
		panic("timeout should be larger than 0")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(c.Addr))

	if err != nil {
		log.Error("open mongodb error(%v)", err)
		panic(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(c.Timeout))
	err = client.Connect(ctx)
	if err != nil {
		log.Error("connect mongodb error(%v)", err)
		panic(err)
	}

	// connect will success even if mongod is closed
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Error("ping mongodb error(%v)", err)
		panic(err)
	}

	log.Info("mongodb connected")
	database := client.Database(c.DBName)

	return &DB{ctx, client, database}
}