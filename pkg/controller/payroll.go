package controller

import (
	"Penggajian/pkg/model"
	"Penggajian/pkg/util"
	"encoding/base64"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

func (a *API) GetPayrollById(c *fiber.Ctx) error {
	// Fetch parameter
	id := c.Params("id")
	if util.IsEmpty(id) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_PARAMETER)
	}

	// Parse into ObjectID
	staffId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_MALFORM)
	}

	payrolls, err := a.payrollRepo.GetPayrollByStaffId(staffId)
	if err != nil {
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

	payrolls, err := a.payrollRepo.GetPayrolls(body.Months, body.Years)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, "cannot fetch payroll contents")
	}

	return SendSuccessResponse(c, fasthttp.StatusOK, payrolls)
}

func (a *API) UpdatePayroll(c *fiber.Ctx) error {
	// Fetch parameter
	id := c.Params("id")
	if util.IsEmpty(id) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_PARAMETER)
	}

	// Parse into ObjectID
	staffId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
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

	payroll, err := a.payrollRepo.EditAndFindById(staffId, body.Months, body.Years, &data)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusNotFound, "no payroll data")
	}

	return SendSuccessResponse(c, fasthttp.StatusOK, payroll)
}

func (a *API) CopyPayroll(c *fiber.Ctx) error {
	return SendSuccessResponse(c, fasthttp.StatusCreated, nil)
}

func (a *API) ImportPayrollFromExcel(c *fiber.Ctx) error {
	// Fetch body
	body := model.PayrollRequest{}
	if c.BodyParser(&body) != nil {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_BODY)
	}

	log.Println("test")
	data := body.Data.(string)
	dest, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "failed to get file")
	}

	log.Println(string(dest))

	return SendSuccessResponse(c, fasthttp.StatusCreated, nil)
}
