package repository

import (
	"Penggajian/dbutil"
	"Penggajian/model"
	"Penggajian/util"
	"errors"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"strconv"
)

type PayrollRepository struct {
	db         *dbutil.Database
	collection *mongo.Collection
}

func NewPayrollRepository(db_ *dbutil.Database, collection_ string) PayrollRepository {
	return PayrollRepository{db: db_, collection: db_.DB.Collection(collection_, options.Collection())}
}

func (p *PayrollRepository) ClearPayroll(months_ uint8, years_ uint16) error {
	ctx, cancel := util.CreateShortTimeoutContext()
	defer cancel()

	// Delete payroll for all staff at that months and years
	result, err := p.collection.DeleteMany(ctx, bson.M{"months": months_, "years": years_})
	if util.IsError(err) {
		return err
	}
	// Delete saving for all staff at that months and years

	log.Println("Deleted ", result.DeletedCount, " data")
	return nil
}

func (p *PayrollRepository) RemovePayrollBySerialNumber(serialNumber_ uint16, months_ uint8, years_ uint16) error {
	ctx, cancel := util.CreateShortTimeoutContext()
	defer cancel()

	return util.GetError(p.collection.DeleteOne(ctx, bson.M{"no_urut_staff": serialNumber_, "months": months_, "years": years_}))
}

func (p *PayrollRepository) ImportFromExcel(reader_ io.Reader, months_ uint8, years_ uint16) (int, error) {
	// Check month and years to prevent overwriting existing data
	if err := p.IsPayrollsEmpty(months_, years_); err != nil {
		return 0, err
	}

	file, err := excelize.OpenReader(reader_, excelize.Options{})
	if util.IsError(err) {
		return 0, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()

	// Get sheets
	sheets := file.GetSheetMap()
	if len(sheets) > 1 {
		return 0, errors.New("sheet should be only has one layer")
	}

	usedSheet := ""
	for index, sheet := range sheets {
		// Set active sheet with named one
		if util.IsEmpty(usedSheet) {
			file.SetActiveSheet(index)
			usedSheet = sheet
			break
		}
	}

	// Get cell value
	rows, err := file.GetRows(usedSheet, excelize.Options{})
	if util.IsError(err) {
		return 0, err
	}

	var payrolls []interface{}
	// starts at row 4
	for i := 4; rows[i] != nil; i++ {
		current := rows[i]
		payroll := model.Payroll{
			StaffName: current[1],
			Institute: current[2],
			Position:  current[3],
			Month:     months_,
			Years:     years_,
		}

		// Serial Number
		serialNumber, err := strconv.ParseUint(current[0], 10, 16)
		if util.IsError(err) {
			return 0, err
		}
		payroll.StaffSerialNumber = uint16(serialNumber)

		// Hours
		// Check empty string
		if util.IsEmpty(current[4]) {
			current[4] = "0"
		}
		// Set Data
		hours, err := strconv.ParseUint(current[4], 10, 8)
		if util.IsError(err) {
			return 0, err
		}
		payroll.Hours = uint8(hours)

		// Staff Salary
		// Check empty string
		if util.IsEmpty(current[5]) {
			current[5] = "0"
		}
		positionSalary, err := strconv.ParseUint(current[5], 10, 64)
		if util.IsError(err) {
			return 0, err
		}
		payroll.Salary.StaffSalary = positionSalary

		// Honorary Salary
		// Check empty string
		if util.IsEmpty(current[6]) {
			current[6] = "0"
		}
		honorarySalary, err := strconv.ParseUint(current[6], 10, 64)
		if util.IsError(err) {
			return 0, err
		}
		payroll.Salary.HonorarySalary = honorarySalary

		// Homeroom Salary
		// Check empty string
		if util.IsEmpty(current[7]) {
			current[7] = "0"
		}
		homeroomSalary, err := strconv.ParseUint(current[7], 10, 64)
		if util.IsError(err) {
			return 0, err
		}
		payroll.Salary.HomeRoomTeacherSalary = homeroomSalary

		// TPMPS
		// Check empty string
		if util.IsEmpty(current[8]) {
			current[8] = "0"
		}
		tpmps, err := strconv.ParseUint(current[8], 10, 64)
		if util.IsError(err) {
			return 0, err
		}
		payroll.Salary.TPMPS = tpmps

		// Saving
		// Check empty string
		if util.IsEmpty(current[9]) {
			current[9] = "0"
		}
		save, err := strconv.ParseUint(current[9], 10, 64)
		if util.IsError(err) {
			return 0, err
		}
		payroll.Save = save

		// Nilam Salary Cut
		// Check empty string
		if util.IsEmpty(current[10]) {
			current[10] = "0"
		}
		nilam, err := strconv.ParseUint(current[10], 10, 64)
		if util.IsError(err) {
			return 0, err
		}
		payroll.SalaryCuts.Nilam = nilam

		// BPJS TK SMP Salary Cut
		// Check empty string
		if util.IsEmpty(current[11]) {
			current[11] = "0"
		}
		bpjsTkSmp, err := strconv.ParseUint(current[11], 10, 64)
		if util.IsError(err) {
			return 0, err
		}
		payroll.SalaryCuts.BPJSTKSMP = bpjsTkSmp

		// BPJS SMP Salary Cut
		// Check empty string
		if util.IsEmpty(current[12]) {
			current[12] = "0"
		}
		bpjsSmp, err := strconv.ParseUint(current[12], 10, 64)
		if util.IsError(err) {
			return 0, err
		}
		payroll.SalaryCuts.BPJSSMP = bpjsSmp

		// BPJS TK SMK Salary Cut
		// Check empty string
		if util.IsEmpty(current[13]) {
			current[13] = "0"
		}
		bpjsTkSmk, err := strconv.ParseUint(current[13], 10, 64)
		if util.IsError(err) {
			return 0, err
		}
		payroll.SalaryCuts.BPJSTKSMK = bpjsTkSmk

		// BPJS SMK Salary Cut
		// Check empty string
		if util.IsEmpty(current[14]) {
			current[14] = "0"
		}
		bpjsSmk, err := strconv.ParseUint(current[14], 10, 64)
		if util.IsError(err) {
			return 0, err
		}
		payroll.SalaryCuts.BPJSSMK = bpjsSmk
		payroll.Total = nil

		payrolls = append(payrolls, payroll)
	}

	// Insert to database
	result, err := p.AddPayrollsInBulk(payrolls)

	return len(result), err
}

func (p *PayrollRepository) GetPayrollsByStaffSerialNumber(serialNum_ uint32, calculateTotal_ bool) ([]model.Payroll, error) {
	return p.getPayrollByFilter(bson.M{"no_urut_staff": serialNum_}, calculateTotal_)
}

func (p *PayrollRepository) GetPayrolls(month_ uint8, years_ uint16, calculateTotal_ bool) ([]model.Payroll, error) {
	payrolls, err := p.getPayrollByFilter(bson.M{"months": month_, "years": years_}, calculateTotal_)
	if util.IsError(err) {
		return payrolls, err
	}

	return payrolls, nil
}

func (p *PayrollRepository) EditAndFindPayrollById(payrollId primitive.ObjectID, months_ uint8, years_ uint16, payroll_ *model.Payroll, calculateTotal_ bool) (model.Payroll, error) {
	err := p.EditPayroll(payrollId, months_, years_, payroll_)
	if err != nil {
		return model.Payroll{}, err
	}

	res, err := p.getPayrollByFilter(bson.M{"_id": payrollId, "months": months_, "years": years_}, calculateTotal_)
	if err != nil {
		return model.Payroll{}, err
	}
	if len(res) < 1 {
		return model.Payroll{}, errors.New("no payroll data")
	}

	return res[0], nil
}

func (p *PayrollRepository) EditPayroll(payrollId primitive.ObjectID, months_ uint8, years_ uint16, payroll_ *model.Payroll) error {
	ctx, cancel := util.CreateShortTimeoutContext()
	defer cancel()

	object, err := util.GenerateBsonObject(*payroll_)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": payrollId, "months": months_, "years": years_}
	query := bson.M{"$set": object}
	return util.GetError(p.collection.UpdateOne(ctx, filter, query, options.Update()))
}

func (p *PayrollRepository) getPayrollByFilter(filter_ any, calculateTotal_ bool) ([]model.Payroll, error) {
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
		if calculateTotal_ {
			payroll.SetTotal()
		}
		payrolls = append(payrolls, payroll)
	}

	return payrolls, nil
}

func (p *PayrollRepository) AddPayrollsInBulk(data_ []any) ([]any, error) {
	ctx, cancel := util.CreateShortTimeoutContext()
	defer cancel()

	queryResult, err := p.collection.InsertMany(ctx, data_, options.InsertMany())
	if util.IsError(err) {
		return nil, err
	}

	return queryResult.InsertedIDs, nil
}

func (p *PayrollRepository) IsPayrollsEmpty(months_ uint8, years_ uint16) error {
	ctx, cancel := util.CreateShortTimeoutContext()
	defer cancel()
	// Check month and years to prevent overwriting existing data
	count, err := p.collection.CountDocuments(ctx, bson.M{"months": months_, "years": years_}, options.Count())
	if util.IsError(err) {
		return err
	}
	if count > 0 {
		return errors.New("cannot overwrite data")
	}

	return nil
}
