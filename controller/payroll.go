package controller

import (
	"Penggajian/model"
	"Penggajian/util"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"strings"
)

func (a *API) GetPayrollBySerialNumber(c *fiber.Ctx) error {
	// Fetch parameter
	id := c.Params("id")
	if util.IsEmpty(id) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_PARAMETER)
	}

	serialNumber, err := strconv.ParseUint(id, 10, 32)
	if util.IsError(err) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_MALFORM)
	}

	payrolls, err := a.payrollRepo.GetPayrollsByStaffSerialNumber(uint32(serialNumber), true)
	if util.IsError(err) {
		return SendErrorResponse(c, fasthttp.StatusNotFound, "no payroll data found")
	}

	return SendSuccessResponse(c, fasthttp.StatusOK, payrolls)
}

func (a *API) GetPayrolls(c *fiber.Ctx) error {

	// Fetch body
	body := model.PayrollRequest{}
	if c.BodyParser(&body) != nil {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_BODY)
	}

	payrolls, err := a.payrollRepo.GetPayrolls(body.Months, body.Years, true)
	if util.IsError(err) {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, "cannot fetch payroll contents")
	}

	return SendSuccessResponse(c, fasthttp.StatusOK, payrolls)
}

func (a *API) UpdatePayrollById(c *fiber.Ctx) error {
	// Fetch parameter
	id := c.Params("id")
	if util.IsEmpty(id) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_PARAMETER)
	}

	// Parse into ObjectID
	payrollId, err := primitive.ObjectIDFromHex(id)
	if util.IsError(err) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_MALFORM)
	}

	// Fetch body
	body := model.PayrollRequest{}
	if c.BodyParser(&body) != nil {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_BODY)
	}

	data, ok := body.Data.(model.Payroll)
	if !ok {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "Malformed Payroll struct")
	}

	payroll, err := a.payrollRepo.EditAndFindPayrollById(payrollId, body.Months, body.Years, &data, true)
	if util.IsError(err) {
		return SendErrorResponse(c, fasthttp.StatusNotFound, "no payroll data")
	}

	return SendSuccessResponse(c, fasthttp.StatusOK, payroll)
}

func (a *API) CopyPayroll(c *fiber.Ctx) error {
	var payrollReq model.PayrollRequest

	if c.BodyParser(&payrollReq) != nil {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_BODY)
	}

	if payrollReq.Months < 1 || payrollReq.Months > 12 {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "bad months")
	}

	// Get last month data
	// Decrement months
	var months uint8 = payrollReq.Months - 1
	var years uint16 = payrollReq.Years
	if payrollReq.Months-1 < 1 {
		months = 12
		years -= 1
	}

	// Insert into new data
	if err := a.payrollRepo.IsPayrollsEmpty(payrollReq.Months, payrollReq.Years); util.IsError(err) {
		return SendErrorResponse(c, fasthttp.StatusConflict, "payroll for current month and years is already there")
	}

	// Get payroll data without calculated total
	payrolls, err := a.payrollRepo.GetPayrolls(months, years, false)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, err.Error())
	}
	// Empty payroll
	if len(payrolls) < 1 {
		return SendErrorResponse(c, fasthttp.StatusNotFound, "no payroll data found in the last month")
	}

	// Increment months, convert into []interface{}, and insert savings
	insertPayrolls := make([]any, len(payrolls))
	for i := 0; i < len(payrolls); i++ {
		// Update month and years based on req (basically just incrementing month by one)
		payrolls[i].Month = payrollReq.Months
		payrolls[i].Years = payrollReq.Years

		// Delete id
		payrolls[i].Id = primitive.NilObjectID

		// Insert array in staff savings
		saving := model.Saving{
			Total:  payrolls[i].Save,
			Months: payrolls[i].Month,
			Years:  payrolls[i].Years,
		}
		if util.IsError(
			util.GetError(a.staffRepo.AddSavingBySerialNumber(payrolls[i].StaffSerialNumber, &saving))) {
			return SendErrorResponse(c, fasthttp.StatusInternalServerError, "cannot update saving data")
		}

		insertPayrolls[i] = payrolls[i]
	}

	size, err := a.payrollRepo.AddPayrollsInBulk(insertPayrolls)
	if util.IsError(err) {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, err.Error())
	}

	return SendSuccessResponse(c, fasthttp.StatusCreated, size)
}

func (a *API) ImportPayrollFromExcel(c *fiber.Ctx) error {
	monthsStr := c.FormValue("months", "")
	yearsStr := c.FormValue("years", "")

	// Check required form value
	if util.IsEmpty(monthsStr) || util.IsEmpty(yearsStr) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_BODY)
	}

	months, err := strconv.ParseUint(monthsStr, 10, 8)
	if util.IsError(err) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "malformed value")
	}
	years, err := strconv.ParseUint(yearsStr, 10, 16)
	if util.IsError(err) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "malformed value")
	}

	// Get form file
	form, err := c.FormFile("data")
	if util.IsError(err) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "data is not satisfied")
	}

	// Check formats (XLAM / XLSM / XLSX / XLTM / XLTX)
	if !strings.Contains(form.Filename, ".xlsx") &&
		!strings.Contains(form.Filename, ".xlam") &&
		!strings.Contains(form.Filename, ".xlsm") &&
		!strings.Contains(form.Filename, ".xltm") &&
		!strings.Contains(form.Filename, ".xltx") {
		return SendErrorResponse(c, fasthttp.StatusConflict, "bad file format")
	}

	files, err := form.Open()
	if util.IsError(err) {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, "failed to open file")
	}

	// Import into payroll repo
	_, err = a.payrollRepo.ImportFromExcel(files, uint8(months), uint16(years))
	if util.IsError(err) {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, "failed to read file")
	}

	// Fetch data from database
	payrolls, err := a.payrollRepo.GetPayrolls(uint8(months), uint16(years), true)
	if util.IsError(err) {
		return SendErrorResponse(c, fasthttp.StatusNotFound, "no payroll data found")
	}

	// Insert array in staff savings
	for i := 0; i < len(payrolls); i++ {
		saving := model.Saving{
			Total:  payrolls[i].Save,
			Months: uint8(months),
			Years:  uint16(years),
		}
		if util.IsError(
			util.GetError(a.staffRepo.AddSavingBySerialNumber(payrolls[i].StaffSerialNumber, &saving))) {
			return SendErrorResponse(c, fasthttp.StatusInternalServerError, "cannot update saving data")
		}
	}

	return SendSuccessResponse(c, fasthttp.StatusCreated, payrolls)
}

func (a *API) ClearPayroll(c *fiber.Ctx) error {
	var payrollReq model.PayrollRequest

	if c.BodyParser(&payrollReq) != nil {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_BODY)
	}

	// Clear payroll
	err := a.payrollRepo.ClearPayroll(payrollReq.Months, payrollReq.Years)
	if util.IsError(err) {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, "failed to clear payroll data")
	}

	// Clear saving
	err = a.staffRepo.ClearSavings(payrollReq.Months, payrollReq.Years)
	if util.IsError(err) {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, "failed to clear saving data")
	}

	return SendSuccessResponse(c, fasthttp.StatusOK, nil)
}
