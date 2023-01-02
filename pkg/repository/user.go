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

type UserRepository struct {
	db         *dbutil.Database
	collection *mongo.Collection
}

func NewUserRepository(db_ *dbutil.Database, collectionName_ string) UserRepository {
	return UserRepository{db: db_, collection: db_.DB.Collection(collectionName_, options.Collection())}
}

func (r *UserRepository) AddUser(user_ *model.User) (primitive.ObjectID, error) {
	ctx, cancel := util.CreateTimeoutContext()
	defer cancel()

	// Give Default Value
	user_.SetDefaultValue(util.GetValue(primitive.ObjectIDFromHex("0")))

	// TODO: Get teacher if there is not, return err
	// TODO: Get user with same teacher objectId, if is there then return err

	// Hash password
	hashedPassword, err := util.Hash(user_.Password)
	if err != nil {
		return util.GetValue(primitive.ObjectIDFromHex("0")), err
	}
	user_.Password = hashedPassword

	res, err := r.collection.InsertOne(ctx, *user_, options.InsertOne())
	if err != nil {
		return util.GetValue(primitive.ObjectIDFromHex("0")), err
	}

	return res.InsertedID.(primitive.ObjectID), err
}

func (r *UserRepository) GetUserByName(username_ string) (model.User, error) {
	user := model.User{}
	ctx, cancel := util.CreateTimeoutContext()
	defer cancel()

	result := r.collection.FindOne(ctx, bson.M{"username": username_}, options.FindOne())
	if result != nil {
		return user, errors.New("user not found")
	}

	if err := result.Decode(&user); err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepository) GetUserById(id_ string) (model.User, error) {
	user := model.User{}
	ctx, cancel := util.CreateTimeoutContext()
	defer cancel()

	id, err := primitive.ObjectIDFromHex(id_)
	if err != nil {
		return user, err
	}

	result := r.collection.FindOne(ctx, bson.M{"_id": id}, options.FindOne())
	if result != nil {
		return user, errors.New("user not found")
	}

	if err := result.Decode(&user); err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepository) GetUsers() ([]model.User, error) {
	var users []model.User
	ctx, cancel := util.CreateTimeoutContext()
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find())
	if err != nil {
		return users, err
	}

	for cursor.Next(ctx) {
		var user model.User
		if err = cursor.Decode(&user); err != nil {
			return users, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) ValidateUser(username_ string, password_ string) (model.User, error) {
	user := model.User{}
	ctx, cancel := util.CreateTimeoutContext()
	defer cancel()

	// Hash password
	hashedPassword, err := util.Hash(password_)
	if err != nil {
		return user, err
	}

	// Get user based on username and password
	result := r.collection.FindOne(ctx, bson.M{"username": username_, "password": hashedPassword}, options.FindOne())
	if result != nil {
		return user, errors.New("there is no data")
	}

	if err = result.Decode(&user); err != nil {
		return user, err
	}

	// Check user condition
	if user.IsLoggedIn {
		return user, errors.New("user already logged in")
	}
	// Update logged in condition
	_, err = r.collection.UpdateByID(ctx, user.Id, bson.M{"is_logged_in": true})
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
