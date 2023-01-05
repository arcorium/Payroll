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
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db         *dbutil.Database
	collection *mongo.Collection
}

func NewUserRepository(db_ *dbutil.Database, collectionName_ string) UserRepository {
	return UserRepository{db: db_, collection: db_.DB.Collection(collectionName_, options.Collection())}
}

func (u *UserRepository) AddUser(user_ *model.User) (primitive.ObjectID, error) {
	ctx, cancel := util.CreateTimeoutContext()
	defer cancel()

	// Give Default Value
	user_.SetDefaultValue()

	// Check existence
	if u.isExist(user_) {
		return primitive.NilObjectID, errors.New("data already exists")
	}

	// Hash password
	hashedPassword, err := util.Hash(user_.Password)
	if err != nil {
		return primitive.NilObjectID, err
	}
	user_.Password = hashedPassword

	res, err := u.collection.InsertOne(ctx, *user_, options.InsertOne())
	if err != nil {
		return primitive.NilObjectID, err
	}

	return res.InsertedID.(primitive.ObjectID), err
}

func (u *UserRepository) RemoveUserById(id_ primitive.ObjectID) (model.User, error) {
	return u.removeUserByFilter(bson.M{"_id": id_})
}

func (u *UserRepository) RemoveUserByName(username_ string) (model.User, error) {
	return u.removeUserByFilter(bson.M{"username": username_})
}

func (u *UserRepository) GetUserByName(username_ string) (model.User, error) {
	return u.getUserByFilter(bson.M{"username": username_})
}

func (u *UserRepository) GetUserById(id_ primitive.ObjectID) (model.User, error) {
	return u.getUserByFilter(bson.M{"_id": id_})
}

func (u *UserRepository) GetUsers() ([]model.User, error) {
	// Create context
	var users []model.User
	ctx, cancel := util.CreateTimeoutContext()
	defer cancel()

	// Find data with no filter
	cursor, err := u.collection.Find(ctx, bson.M{}, options.Find())
	if err != nil {
		return users, err
	}

	// Iterate cursor
	for cursor.Next(ctx) {
		var user model.User
		if err = cursor.Decode(&user); err != nil {
			return users, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (u *UserRepository) EditUserById(id_ primitive.ObjectID, user_ *model.User) (primitive.ObjectID, error) {
	return u.editUserByFilter(bson.M{"_id": id_}, user_)
}

func (u *UserRepository) EditUserByName(username_ string, user_ *model.User) (primitive.ObjectID, error) {
	return u.editUserByFilter(bson.M{"username": username_}, user_)
}

func (u *UserRepository) ValidateUser(username_ string, password_ string) (model.User, error) {
	// Get user based on username
	user, err := u.GetUserByName(username_)
	if err != nil {
		return user, err
	}
	// Matching password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password_))

	return user, err
}

func (u *UserRepository) UpdateLoggedIn(userId_ primitive.ObjectID, condition_ bool) error {
	ctx, cancel := util.CreateTimeoutContext()
	defer cancel()

	// Update logged in condition
	return util.GetError(u.collection.UpdateByID(ctx, userId_, bson.M{"$set": bson.M{"is_logged_in": condition_}}))
}

func (u *UserRepository) IsLoggedIn(userId_ primitive.ObjectID) error {
	user, err := u.GetUserById(userId_)
	if err != nil {
		return err
	}

	if !user.IsLoggedIn {
		return errors.New("user not logged in")
	}

	return nil
}

func (u *UserRepository) getUserByFilter(filter_ any) (model.User, error) {
	user := model.User{}
	// Create Timeout context
	ctx, cancel := util.CreateTimeoutContext()
	defer cancel()

	// Get
	result := u.collection.FindOne(ctx, filter_, options.FindOne())
	if result == nil {
		return user, errors.New("user not found")
	}

	// Decode
	err := result.Decode(&user)
	return user, err
}

func (u *UserRepository) removeUserByFilter(filter_ any) (model.User, error) {
	user := model.User{}
	// Create Context
	ctx, cancel := util.CreateTimeoutContext()
	defer cancel()

	// Find and Delete
	result := u.collection.FindOneAndDelete(ctx, filter_, options.FindOneAndDelete())
	if result == nil {
		return user, errors.New("user not found")
	}

	// Decode
	err := result.Decode(&user)
	return user, err
}

func (u *UserRepository) editUserByFilter(filter_ any, user_ *model.User) (primitive.ObjectID, error) {
	// Create context
	ctx, cancel := util.CreateTimeoutContext()
	defer cancel()

	// Update modified_at field
	user_.UpdateModifiedTime()

	// Replace
	result, err := u.collection.ReplaceOne(ctx, filter_, *user_, options.Replace())
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.UpsertedID.(primitive.ObjectID), nil
}

func (u *UserRepository) isExist(user_ *model.User) bool {
	return util.GetError(u.GetUserByName(user_.Username)) == nil
}
