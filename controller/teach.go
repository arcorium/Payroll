package controller

import (
	"Penggajian/model"
	"Penggajian/util"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

func (a *API) UpdateTeachTime(c *fiber.Ctx) error {
	id := c.Params("id")
	if util.IsEmpty(id) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_PARAMETER)
	}
	uuid := c.Params("uuid")
	if util.IsEmpty(uuid) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_PARAMETER)
	}

	// Parse body
	teachDetails := model.TeachTimeDetail{}
	if c.BodyParser(&teachDetails) != nil {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_BODY)
	}

	// Parse id
	staffId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_MALFORM)
	}

	staff, err := a.staffRepo.EditTeachTime(staffId, uuid, &teachDetails)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, "failed to update staff")
	}

	return SendSuccessResponse(c, fasthttp.StatusOK, staff)
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
