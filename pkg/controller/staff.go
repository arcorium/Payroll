package controller

import (
	"Penggajian/pkg/model"
	"Penggajian/pkg/util"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *API) RegisterStaff(c *fiber.Ctx) error {
	staff := model.Staff{}

	// Parse body
	if err := c.BodyParser(&staff); err != nil {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_BODY)
	}

	// Insert data
	response, err := a.staffRepo.AddStaff(&staff)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, "failed to add staff")
	}

	// Get inserted data
	staff, err = a.staffRepo.GetStaffById(response)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusNotFound, "there is no staff data")
	}

	// Response
	return SendSuccessResponse(c, fasthttp.StatusCreated, staff)
}

func (a *API) UpdateStaffById(c *fiber.Ctx) error {
	staff := model.Staff{}
	if c.BodyParser(&staff) != nil {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_BODY)
	}

	// Pass data into repository
	updatedStaff, err := a.staffRepo.EditAndFindById(staff.Id, &staff)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusNotFound, "failed to set staff status")
	}

	return SendSuccessResponse(c, fasthttp.StatusOK, updatedStaff)
}

func (a *API) GetStaffById(c *fiber.Ctx) error {
	// Parameter check
	id := c.Params("id")
	if util.IsEmpty(id) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_PARAMETER)
	}

	// Convert into ObjectID
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Pass data into repository
	staff, err := a.staffRepo.GetStaffById(objectId)
	if err != nil {
		return err
	}

	// Response
	return SendSuccessResponse(c, fasthttp.StatusOK, staff)
}

func (a *API) GetStaffs(c *fiber.Ctx) error {
	// Get data from repository
	staffs, err := a.staffRepo.GetStaffs()
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, "failed to retrieve staffs data")
	}

	// Response
	return SendSuccessResponse(c, fasthttp.StatusOK, staffs)
}

func (a *API) InsertTeachTime(c *fiber.Ctx) error {
	id := c.Params("id")
	if util.IsEmpty(id) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_PARAMETER)
	}

	// Parse id
	staffId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_MALFORM)
	}

	// Parse body
	teachDetails := model.TeachTimeDetail{}
	if c.BodyParser(&teachDetails) != nil {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_BODY)
	}

	// Insert data
	_, err = a.staffRepo.AddTeachTime(staffId, &teachDetails)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, "failed to add teach time")
	}

	return SendSuccessResponse(c, fasthttp.StatusCreated, teachDetails)
}

func (a *API) RemoveTeachTime(c *fiber.Ctx) error {
	id := c.Params("id")
	if util.IsEmpty(id) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_PARAMETER)
	}
	uuid := c.Params("uuid")
	if util.IsEmpty(uuid) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_PARAMETER)
	}

	// Parse id
	staffId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_MALFORM)
	}

	err = a.staffRepo.RemoveTeachTime(staffId, uuid)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, "failed to remove teach time")
	}

	return SendSuccessResponse(c, fasthttp.StatusOK, nil)
}
