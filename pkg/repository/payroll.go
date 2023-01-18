package repository

import (
	"Penggajian/pkg/dbutil"
	"Penggajian/pkg/model"
	"Penggajian/pkg/util"
	"errors"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
)

type PayrollRepository struct {
	db         *dbutil.Database
	collection *mongo.Collection
}

func NewPayrollRepository(db_ *dbutil.Database, collection_ string) PayrollRepository {
	return PayrollRepository{db: db_, collection: db_.DB.Collection(collection_, options.Collection())}
}

func (p *PayrollRepository) ImportFromExcel(reader_ io.Reader, years_ uint16) error {
	ctx, cancel := util.CreateShortTimeoutContext()
	defer cancel()

	file, err := excelize.OpenReader(reader_, excelize.Options{})
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()

	// Get cell value

	// Set into database

	// Drop from database
	result, err := p.collection.DeleteMany(ctx, bson.M{"years": years_}, options.Delete())
	if err != nil {
		return err
	}
	log.Println("Deleted Payroll in years ", years_, " count ", result.DeletedCount)

	// Insert to database

	return nil
}

func (p *PayrollRepository) GetPayrollByStaffId(staffId_ primitive.ObjectID) ([]model.Payroll, error) {
	return p.getPayrollByFilter(bson.M{"staff_id": staffId_})
}

func (p *PayrollRepository) GetOnePayrollByStaffId(staffId_ primitive.ObjectID, months_ uint8, years_ uint16) (model.Payroll, error) {
	payrolls, err := p.getPayrollByFilter(bson.M{"staff_id": staffId_, "months": months_, "years": years_})
	if err != nil {
		return model.Payroll{}, err
	}

	return payrolls[0], nil
}

func (p *PayrollRepository) GetPayrolls(month_ uint8, years_ uint16) ([]model.Payroll, error) {
	return p.getPayrollByFilter(bson.M{"months": month_, "years": years_})
}

func (p *PayrollRepository) EditAndFindById(staffId_ primitive.ObjectID, months_ uint8, years_ uint16, payroll_ *model.Payroll) (model.Payroll, error) {
	err := p.EditPayroll(staffId_, months_, years_, payroll_)
	if err != nil {
		return model.Payroll{}, err
	}

	return p.GetOnePayrollByStaffId(staffId_, months_, years_)
}

func (p *PayrollRepository) EditPayroll(staffId_ primitive.ObjectID, months_ uint8, years_ uint16, payroll_ *model.Payroll) error {
	ctx, cancel := util.CreateShortTimeoutContext()
	defer cancel()

	object, err := util.GenerateBsonObject(*payroll_)
	if err != nil {
		return err
	}

	filter := bson.M{"staff_id": staffId_, "months": months_, "years": years_}
	query := bson.M{"$set": object}
	return util.GetError(p.collection.UpdateOne(ctx, filter, query, options.Update()))
}

func (p *PayrollRepository) getPayrollByFilter(filter_ any) ([]model.Payroll, error) {
	var payrolls []model.Payroll
	ctx, cancel := util.CreateShortTimeoutContext()
	defer cancel()

	cursor, err := p.collection.Find(ctx, filter_, options.Find())
	if err != nil {
		return payrolls, errors.New("there is no payroll data")
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		payroll := model.Payroll{}
		err = cursor.Decode(&payroll)
		if err != nil {
			return payrolls, err
		}
		payrolls = append(payrolls, payroll)
	}

	return payrolls, nil
}
