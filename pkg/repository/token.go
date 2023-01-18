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
	ctx, cancel := util.CreateShortTimeoutContext()
	defer cancel()

	result, err := t.collection.InsertOne(ctx, *token_, options.InsertOne())
	if err != nil {
		return model.NullResponseID(), err
	}

	return model.NewResponseID(result.InsertedID.(primitive.ObjectID)), nil
}

// UpsertTokenByUserId returned token will be nil if it inserting new data
func (t *TokenRepository) UpsertTokenByUserId(userId_ primitive.ObjectID, token_ *model.Token) (model.Token, error) {
	return t.upsertTokenByFilter(bson.M{"user_id": userId_}, bson.M{"refresh_token": token_.Token})
}

func (t *TokenRepository) RemoveTokenByToken(token_ string) (model.Token, error) {
	return t.removeTokenByFilter(bson.M{"refresh_token": token_})
}

func (t *TokenRepository) RemoveTokenById(id_ primitive.ObjectID) (model.Token, error) {
	return t.removeTokenByFilter(bson.M{"_id": id_})
}

func (t *TokenRepository) RemoveTokenByUserId(userId_ primitive.ObjectID) (model.Token, error) {
	return t.removeTokenByFilter(bson.M{"user_id": userId_})
}

func (t *TokenRepository) UpdateToken(token_ string, newToken_ string) error {
	// Create timeout context
	ctx, cancel := util.CreateShortTimeoutContext()
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

func (t *TokenRepository) GetTokenByUserId(userId_ primitive.ObjectID) (model.Token, error) {
	return t.getTokenByFilter(bson.M{"user_id": userId_})
}

func (t *TokenRepository) GetTokenByToken(token_ string) (model.Token, error) {
	return t.getTokenByFilter(bson.M{"refresh_token": token_})
}

func (t *TokenRepository) getTokenByFilter(filter_ any) (model.Token, error) {
	token := model.Token{}

	// Create timeout context
	ctx, cancel := util.CreateShortTimeoutContext()
	defer cancel()

	// Querying
	result := t.collection.FindOne(ctx, filter_, options.FindOne())

	// Decode
	err := result.Decode(&token)
	return token, err
}

func (t *TokenRepository) removeTokenByFilter(filter_ any) (model.Token, error) {
	token := model.Token{}

	// Create timeout context
	ctx, cancel := util.CreateShortTimeoutContext()
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

func (t *TokenRepository) upsertTokenByFilter(filter_ any, update_ any) (model.Token, error) {
	ctx, cancel := util.CreateShortTimeoutContext()
	defer cancel()

	// When filter is not found, it will insert new data with filter field and field after $set
	token := model.Token{}
	result := t.collection.FindOneAndUpdate(ctx, filter_, bson.M{"$set": update_}, options.FindOneAndUpdate().SetUpsert(true))

	err := result.Decode(&token)
	return token, err
}
