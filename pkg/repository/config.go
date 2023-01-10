package repository

import (
	"Penggajian/pkg/dbutil"
	"Penggajian/pkg/model"
	"Penggajian/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConfigRepository struct {
	db         *dbutil.Database
	collection *mongo.Collection
}

func NewConfigRepository(db_ *dbutil.Database, collectionName_ string) ConfigRepository {
	return ConfigRepository{db: db_, collection: db_.DB.Collection(collectionName_, options.Collection())}
}

func (c *ConfigRepository) ReplaceConfig(config_ *model.Config) error {
	ctx, cancel := util.CreateTimeoutContext()
	defer cancel()

	id, err := primitive.ObjectIDFromHex("63bcfde6176616e32e0e2c68")
	if err != nil {
		return err
	}

	_, err = c.collection.ReplaceOne(ctx, bson.M{"_id": id}, *config_, options.Replace())
	if err != nil {
		return err
	}

	return nil
}
