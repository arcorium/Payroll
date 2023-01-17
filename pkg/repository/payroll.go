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
	ctx, cancel := util.CreateTimeoutContext()
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

func (p *PayrollRepository) GetPayrollByStaffId(staffId_ primitive.ObjectID) (model.Payroll, error) {
	return p.getPayrollByFilter(bson.M{"staff_id": staffId_})
}

func (p *PayrollRepository) GetPayrolls() ([]model.Payroll, error) {
	return nil, nil
}

func (p *PayrollRepository) getPayrollByFilter(filter_ any) (model.Payroll, error) {
	ctx, cancel := util.CreateTimeoutContext()
	defer cancel()

	payroll := model.Payroll{}

	res := p.collection.FindOne(ctx, filter_, options.FindOne())
	if res == nil {
		return payroll, errors.New("there is no payroll data")
	}

	err := res.Decode(&payroll)
	return payroll, err
}
