package repository

import (
	"Penggajian/pkg/dbutil"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TeachRepository struct {
	db         *dbutil.Database
	collection *mongo.Collection
}

func NewTeachRepository(db_ *dbutil.Database, collectionName_ string) TeachRepository {
	return TeachRepository{db: db_, collection: db_.DB.Collection(collectionName_, options.Collection())}
}

