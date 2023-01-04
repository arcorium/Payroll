package repository

import (
	"Penggajian/pkg/dbutil"
	"Penggajian/pkg/model"
	"Penggajian/pkg/util"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TokenRepository struct {
	db         *dbutil.Database
	collection *mongo.Collection
}

func NewTokenRepository(db_ *dbutil.Database, collectionName_ string) TokenRepository {
	return TokenRepository{db: db_, collection: db_.DB.Collection(collectionName_)}
}

func (t *TokenRepository) AddToken(token_ *model.Token) (model.ResponseID, error) {
	ctx, cancel := util.CreateTimeoutContext()
	defer cancel()

	result, err := t.collection.InsertOne(ctx, *token_, options.InsertOne())
	if err != nil {
		return model.NullResponseID(), err
	}

	return model.NewResponseID(result.InsertedID.(primitive.ObjectID)), nil
}

func (t *TokenRepository) RemoveTokenByToken(token_ string) (model.Token, error) {
	return t.removeTokenByFilter(bson.M{"refresh_token": token_})
}

func (t *TokenRepository) RemoveTokenById(id_ primitive.ObjectID) (model.Token, error) {
	return t.removeTokenByFilter(bson.M{"_id": id_})
}

func (t *TokenRepository) UpdateToken(token_ string, newToken_ string) error {
	// Create timeout context
	ctx, cancel := util.CreateTimeoutContext()
	defer cancel()

	// Update query
	result, err := t.collection.UpdateOne(ctx, bson.M{"refresh_token": token_}, bson.M{"$set": bson.M{"refresh_token": newToken_}}, options.Update())
	if err != nil {
		return err
	}

	// Check total modified
	if result.ModifiedCount < 1 {
		return errors.New("token not found")
	}

	return nil
}

func (t *TokenRepository) ValidateToken(token_ string) (model.Token, error) {
	return t.getToken(token_)
}

func (t *TokenRepository) getToken(token_ string) (model.Token, error) {
	token := model.Token{}

	// Create timeout context
	ctx, cancel := util.CreateTimeoutContext()
	defer cancel()

	// Querying
	result := t.collection.FindOne(ctx, bson.M{"refresh_token": token_}, options.FindOne())
	if result == nil {
		return token, errors.New("token not found")
	}

	// Decode
	err := result.Decode(&token)
	return token, err
}

func (t *TokenRepository) removeTokenByFilter(filter_ any) (model.Token, error) {
	token := model.Token{}

	// Create timeout context
	ctx, cancel := util.CreateTimeoutContext()
	defer cancel()

	// Delete and return deleted one
	result := t.collection.FindOneAndDelete(ctx, filter_, options.FindOneAndDelete())
	if result == nil {
		return token, errors.New("token not found")
	}

	// Decode into token structure
	err := result.Decode(&token)
	return token, err
}
