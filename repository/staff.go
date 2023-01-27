package repository

import (
	"Penggajian/dbutil"
	"Penggajian/model"
	"Penggajian/util"
	"errors"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type StaffRepository struct {
	db         *dbutil.Database
	collection *mongo.Collection
}

func NewStaffRepository(db_ *dbutil.Database, collectionName_ string) StaffRepository {
	return StaffRepository{db: db_, collection: db_.DB.Collection(collectionName_, options.Collection())}
}

func (t *StaffRepository) AddStaff(staff_ *model.Staff) (primitive.ObjectID, error) {
	ctx, cancel := util.CreateShortTimeoutContext()
	defer cancel()

	if len(staff_.TeachTimeDetails) < 1 {
		staff_.TeachTimeDetails = []model.TeachTimeDetail{}
	}

	// Generate UUID for each teach time details
	for i := 0; i < len(staff_.TeachTimeDetails); i++ {
		staff_.TeachTimeDetails[i].UUID = uuid.NewString()
	}

	if len(staff_.Savings) < 1 {
		staff_.Savings = []model.Saving{}
	}

	// Generate UUID for each saving
	for i := 0; i < len(staff_.Savings); i++ {
		staff_.Savings[i].UUID = uuid.NewString()
	}

	result, err := t.collection.InsertOne(ctx, *staff_, options.InsertOne())
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (t *StaffRepository) EditById(staffId_ primitive.ObjectID, staff_ *model.Staff) error {
	ctx, cancel := util.CreateShortTimeoutContext()
	defer cancel()

	update, err := util.GenerateBsonObject(*staff_)

	res, err := t.collection.UpdateByID(ctx, staffId_, bson.M{"$set": update}, options.Update())
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.New("no data updated")
	}

	return nil
}

func (t *StaffRepository) EditAndFindById(staffId_ primitive.ObjectID, staff_ *model.Staff) (model.Staff, error) {
	staff := model.Staff{}

	err := t.EditById(staffId_, staff_)
	if err != nil {
		return staff, err
	}

	staff, err = t.GetStaffById(staffId_)
	return staff, err
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
	ctx, cancel := util.CreateShortTimeoutContext()
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
	ctx, cancel := util.CreateShortTimeoutContext()
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
	ctx, cancel := util.CreateShortTimeoutContext()
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
	ctx, cancel := util.CreateShortTimeoutContext()
	defer cancel()

	// Replace
	result, err := t.collection.ReplaceOne(ctx, filter_, *teacher_, options.Replace())
	if err != nil {
		return model.NullResponseID(), err
	}

	return model.NewResponseID(result.UpsertedID.(primitive.ObjectID)), nil
}

// Add Teach Time Details
func (t *StaffRepository) AddTeachTime(staffId_ primitive.ObjectID, details_ *model.TeachTimeDetail) (string, error) {
	// Set default value
	if details_.Months < 1 || details_.Years < 1 {
		year, month, _ := time.Now().Date()
		details_.UUID = uuid.NewString()
		details_.Months = uint8(month)
		details_.Years = uint16(year)
	}

	// Push array
	update := bson.M{"details": *details_}
	err := t.pushArray(bson.M{"_id": staffId_}, update)

	return details_.UUID, err
}

func (t *StaffRepository) RemoveTeachTime(staffId_ primitive.ObjectID, uuid_ string) error {
	ctx, cancel := util.CreateShortTimeoutContext()
	defer cancel()

	update := bson.M{"$pull": bson.M{"details": bson.M{"uuid": uuid_}}}

	return util.GetError(t.collection.UpdateByID(ctx, staffId_, update, options.Update()))
}

func (t *StaffRepository) ClearTeachTime(staffId_ primitive.ObjectID, months_ uint8, years_ uint16) error {
	ctx, cancel := util.CreateLongTimeoutContext()
	defer cancel()

	filter := bson.M{"_id": staffId_}
	query := bson.M{"$pull": bson.M{"details": bson.M{"month": months_, "years": years_}}}

	return util.GetError(t.collection.UpdateOne(ctx, filter, query, options.Update()))
}

func (t *StaffRepository) AddSavingBySerialNumber(serialNumber_ uint16, details_ *model.Saving) (string, error) {
	// Set default value
	details_.UUID = uuid.NewString()
	if details_.Months < 1 || details_.Years < 1 {
		year, month, _ := time.Now().Date()
		details_.Months = uint8(month)
		details_.Years = uint16(year)
	}

	err := t.pushArray(bson.M{"no_urut": serialNumber_}, bson.M{"tabungan": *details_})
	return details_.UUID, err
}

func (t *StaffRepository) RemoveStaffSavingById(staffId_ primitive.ObjectID, uuid_ string) error {
	ctx, cancel := util.CreateShortTimeoutContext()
	defer cancel()

	update := bson.M{"$pull": bson.M{"tabungan": bson.M{"uuid": uuid_}}}

	return util.GetError(t.collection.UpdateByID(ctx, staffId_, update, options.Update()))
}

func (t *StaffRepository) ClearStaffSavingsByYears(staffId_ primitive.ObjectID, years_ uint16) error {
	ctx, cancel := util.CreateLongTimeoutContext()
	defer cancel()

	filter := bson.M{"_id": staffId_}
	query := bson.M{"$pull": bson.M{"tabungan": bson.M{"years": years_}}}

	return util.GetError(t.collection.UpdateOne(ctx, filter, query, options.Update()))
}

func (t *StaffRepository) ClearSavings(months_ uint8, years_ uint16) error {
	ctx, cancel := util.CreateLongTimeoutContext()
	defer cancel()

	filter := bson.M{}
	query := bson.M{"$pull": bson.M{"tabungan": bson.M{"months": months_, "years": years_}}}

	return util.GetError(t.collection.UpdateMany(ctx, filter, query, options.Update()))
}

func (t *StaffRepository) pushArray(filter_ any, query_ any) error {
	ctx, cancel := util.CreateShortTimeoutContext()
	defer cancel()

	query := bson.M{"$push": query_}
	return util.GetError(t.collection.UpdateOne(ctx, filter_, query, options.Update()))
}
