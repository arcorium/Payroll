package dbutil

import (
	"Penggajian/pkg/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	client *mongo.Client
	DB     *mongo.Database
	Config *DBConfig
}

func (d *Database) Disconnect() error {
	ctx, cancel := util.CreateTimeoutContext()
	defer cancel()

	return d.client.Disconnect(ctx)
}

func Connect(config_ *DBConfig) (Database, error) {
	db := Database{Config: config_}
	ctx, cancel := util.CreateTimeoutContext()
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config_.DatabaseURI))
	if err != nil {
		return db, err
	}

	db.client = client
	db.DB = db.client.Database(config_.DatabaseName, options.Database())

	return db, nil
}
