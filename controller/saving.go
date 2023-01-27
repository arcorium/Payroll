package controller

import (
	"Penggajian/model"
	"Penggajian/util"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *API) InsertSaving(c *fiber.Ctx) error {
	id := c.Params("id")
	if util.IsEmpty(id) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_PARAMETER)
	}
	// Parse body
	saving := model.Saving{}
	if util.IsError(c.BodyParser(&saving)) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_BODY)
	}

	// Parse id
	staffId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_MALFORM)
	}

	// Insert
	_, err = a.staffRepo.AddSavingByStaffId(staffId, &saving)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, "failed to add saving")
	}

	return SendSuccessResponse(c, fasthttp.StatusCreated, saving)
}

func (a *API) UpdateSaving(c *fiber.Ctx) error {
	id := c.Params("id")
	if util.IsEmpty(id) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_PARAMETER)
	}
	uuid := c.Params("uuid")
	if util.IsEmpty(uuid) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_PARAMETER)
	}

	// Parse body
	saving := model.Saving{}
	if c.BodyParser(&saving) != nil {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_BODY)
	}

	// Parse id
	staffId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_MALFORM)
	}

	staff, err := a.staffRepo.EditStaffSavingById(staffId, uuid, &saving)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, "failed to update staff")
	}

	return SendSuccessResponse(c, fasthttp.StatusOK, staff)
}

func (a *API) RemoveSaving(c *fiber.Ctx) error {
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

	err = a.staffRepo.RemoveStaffSavingById(staffId, uuid)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, "failed to remove teach time")
	}

	return SendSuccessResponse(c, fasthttp.StatusOK, nil)
}
