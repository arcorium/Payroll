package repository

import (
	"Penggajian/dbutil"
	"Penggajian/model"
	"Penggajian/util"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PositionRepository struct {
	db         *dbutil.Database
	collection *mongo.Collection
}

func NewPositionRepository(db_ *dbutil.Database, collectionName_ string) PositionRepository {
	return PositionRepository{db: db_, collection: db_.DB.Collection(collectionName_)}
}

func (p *PositionRepository) AddPosition(position_ model.Position) (model.ResponseID, error) {
	ctx, cancel := util.CreateShortTimeoutContext()
	defer cancel()

	result, err := p.collection.InsertOne(ctx, position_, options.InsertOne())
	if err != nil {
		return model.NullResponseID(), err
	}

	return model.NewResponseID(result.InsertedID.(primitive.ObjectID)), nil
}

func (p *PositionRepository) RemovePositionByName(name_ string) (model.Position, error) {
	return p.removePositionByFilter(bson.M{"name": name_})
}

func (p *PositionRepository) RemovePositionById(id_ primitive.ObjectID) (model.Position, error) {
	return p.removePositionByFilter(bson.M{"_id": id_})
}

func (p *PositionRepository) GetPositionByName(name_ string) (model.Position, error) {
	return p.getPositionByFilter(bson.M{"name": name_})
}

func (p *PositionRepository) GetPositionById(id_ primitive.ObjectID) (model.Position, error) {
	return p.getPositionByFilter(bson.M{"_id": id_})
}

func (p *PositionRepository) GetPositions() ([]model.Position, error) {
	// Create context
	var positions []model.Position
	ctx, cancel := util.CreateShortTimeoutContext()
	defer cancel()

	// Find data with no filter
	cursor, err := p.collection.Find(ctx, bson.M{}, options.Find())
	if err != nil {
		return positions, err
	}

	// Iterate cursor
	for cursor.Next(ctx) {
		var position model.Position
		if err = cursor.Decode(&position); err != nil {
			return positions, err
		}
		positions = append(positions, position)
	}

	return positions, nil
}

func (p *PositionRepository) EditPositionByName(name_ string, newName_ string) (model.ResponseID, error) {
	return p.editPositionByFilter(bson.M{"name": name_}, newName_)
}

func (p *PositionRepository) EditPositionById(id_ primitive.ObjectID, newName_ string) (model.ResponseID, error) {
	return p.editPositionByFilter(bson.M{"_id": id_}, newName_)
}

func (p *PositionRepository) getPositionByFilter(filter_ any) (model.Position, error) {
	position := model.Position{}
	// Create Timeout context
	ctx, cancel := util.CreateShortTimeoutContext()
	defer cancel()

	// Get
	result := p.collection.FindOne(ctx, filter_, options.FindOne())
	if result == nil {
		return position, errors.New("user not found")
	}

	// Decode
	err := result.Decode(&position)
	return position, err
}

func (p *PositionRepository) removePositionByFilter(filter_ any) (model.Position, error) {
	position := model.Position{}
	// Create Context
	ctx, cancel := util.CreateShortTimeoutContext()
	defer cancel()

	// Find and Delete
	result := p.collection.FindOneAndDelete(ctx, filter_, options.FindOneAndDelete())
	if result == nil {
		return position, errors.New("position not found")
	}

	// Decode
	err := result.Decode(&position)
	return position, err
}

func (p *PositionRepository) editPositionByFilter(filter_ any, newName_ string) (model.ResponseID, error) {
	// Create context
	ctx, cancel := util.CreateShortTimeoutContext()
	defer cancel()

	// Replace
	result, err := p.collection.UpdateOne(ctx, filter_, bson.M{"$set": bson.M{"name": newName_}})
	if err != nil {
		return model.NullResponseID(), err
	}

	return model.NewResponseID(result.UpsertedID.(primitive.ObjectID)), nil
}
