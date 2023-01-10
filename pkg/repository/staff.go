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

type StaffRepository struct {
	db         *dbutil.Database
	collection *mongo.Collection
}

func NewStaffRepository(db_ *dbutil.Database, collectionName_ string) StaffRepository {
	return StaffRepository{db: db_, collection: db_.DB.Collection(collectionName_, options.Collection())}
}

func (t *StaffRepository) AddStaff(teacher_ *model.Staff) (model.ResponseID, error) {
	ctx, cancel := util.CreateTimeoutContext()
	defer cancel()

	result, err := t.collection.InsertOne(ctx, *teacher_, options.InsertOne())
	if err != nil {
		return model.NullResponseID(), err
	}

	return model.NewResponseID(result.InsertedID.(primitive.ObjectID)), nil
}

func (t *StaffRepository) RemoveStaffById(id_ primitive.ObjectID) (model.Staff, error) {
	return t.removeStaffByFilter(bson.M{"_id": id_})
}

func (t *StaffRepository) RemoveStaffByName(name_ string) (model.Staff, error) {
	return t.removeStaffByFilter(bson.M{"name": name_})
}

func (t *StaffRepository) GetStaffByName(name_ string) (model.Staff, error) {
	return t.getStaffByFilter(bson.M{"name:": name_})
}

func (t *StaffRepository) GetStaffById(id_ primitive.ObjectID) (model.Staff, error) {
	return t.getStaffByFilter(bson.M{"_id": id_})
}

func (t *StaffRepository) GetStaffs() ([]model.Staff, error) {
	// Create context
	var teachers []model.Staff
	ctx, cancel := util.CreateTimeoutContext()
	defer cancel()

	// Find data with no filter
	cursor, err := t.collection.Find(ctx, bson.M{}, options.Find())
	if err != nil {
		return teachers, err
	}

	// Iterate cursor
	for cursor.Next(ctx) {
		var teacher model.Staff
		if err = cursor.Decode(&teacher); err != nil {
			return teachers, err
		}
		teachers = append(teachers, teacher)
	}

	return teachers, nil
}

func (t *StaffRepository) EditStaffById(id_ primitive.ObjectID, teacher_ *model.Staff) (model.ResponseID, error) {
	return t.editStaffByFilter(bson.M{"_id": id_}, teacher_)
}

func (t *StaffRepository) EditStaffByName(name_ string, teacher_ *model.Staff) (model.ResponseID, error) {
	return t.editStaffByFilter(bson.M{"name": name_}, teacher_)
}

func (t *StaffRepository) getStaffByFilter(filter_ any) (model.Staff, error) {
	teacher := model.Staff{}
	// Create Timeout context
	ctx, cancel := util.CreateTimeoutContext()
	defer cancel()

	// Get
	result := t.collection.FindOne(ctx, filter_, options.FindOne())
	if result == nil {
		return teacher, errors.New("teacher data not found")
	}

	// Decode
	err := result.Decode(&teacher)
	return teacher, err
}

func (t *StaffRepository) removeStaffByFilter(filter_ any) (model.Staff, error) {
	teacher := model.Staff{}
	// Create Context
	ctx, cancel := util.CreateTimeoutContext()
	defer cancel()

	// Find and Delete
	result := t.collection.FindOneAndDelete(ctx, filter_, options.FindOneAndDelete())
	if result == nil {
		return teacher, errors.New("teacher data not found")
	}

	// Decode
	err := result.Decode(&teacher)
	return teacher, err
}

func (t *StaffRepository) editStaffByFilter(filter_ any, teacher_ *model.Staff) (model.ResponseID, error) {
	// Create context
	ctx, cancel := util.CreateTimeoutContext()
	defer cancel()

	// Replace
	result, err := t.collection.ReplaceOne(ctx, filter_, *teacher_, options.Replace())
	if err != nil {
		return model.NullResponseID(), err
	}

	return model.NewResponseID(result.UpsertedID.(primitive.ObjectID)), nil
}

// Add Teach Time Details
func (t *StaffRepository) AddTeachTime(staffId_ primitive.ObjectID, details_ *model.TeachTimeDetail) (primitive.ObjectID, error) {
	ctx, cancel := util.CreateTimeoutContext()
	defer cancel()

	// Push array
	update := bson.M{"$push": *details_}
	result, err := t.collection.UpdateByID(ctx, staffId_, update, options.Update())
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.UpsertedID.(primitive.ObjectID), nil
}
