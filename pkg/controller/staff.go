package controller

import (
	"Penggajian/pkg/model"
	"Penggajian/pkg/util"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *API) RegisterStaff(c *fiber.Ctx) error {
	teacher := model.Staff{}

	// Parse body
	if err := c.BodyParser(&teacher); err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, "failed to parse body")
	}

	// Response
	response, err := a.staffRepo.AddStaff(&teacher)
	if err != nil {
		return err
	}

	// Response
	err = c.JSON(response)
	return err
}

func (a *API) RemoveStaffByName(c *fiber.Ctx) error {
	// Parameter check
	name := c.Params("name")
	if len(name) < 1 {
		return errors.New("parameter is not satisfied")
	}

	// Pass data into repository
	teacher, err := a.staffRepo.RemoveStaffByName(name)
	if err != nil {
		return err
	}

	// Response
	err = c.JSON(teacher)
	return err
}

func (a *API) RemoveStaffById(c *fiber.Ctx) error {
	// Parameter check
	id := c.Params("name")
	if len(id) < 1 {
		return errors.New("parameter is not satisfied")
	}

	// Convert into ObjectID
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Pass data into repository
	teacher, err := a.staffRepo.RemoveStaffById(objectId)
	if err != nil {
		return err
	}

	// Response
	err = c.JSON(teacher)
	return err
}

func (a *API) GetStaffByName(c *fiber.Ctx) error {
	// Parameter check
	name := c.Params("name")
	if len(name) < 1 {
		return errors.New("parameter is not satisfied")
	}

	// Pass data into repository
	user, err := a.staffRepo.GetStaffByName(name)
	if err != nil {
		return err
	}

	// Response
	err = c.JSON(user)
	return err
}

func (a *API) GetStaffById(c *fiber.Ctx) error {
	// Parameter check
	id := c.Params("id")
	if len(id) < 1 {
		return errors.New("parameter is not satisfied")
	}

	// Convert into ObjectID
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Pass data into repository
	position, err := a.staffRepo.GetStaffById(objectId)
	if err != nil {
		return err
	}

	// Response
	err = c.JSON(position)
	return err
}

func (a *API) GetStaffs(c *fiber.Ctx) error {
	// Get data from repository
	teachers, err := a.staffRepo.GetStaffs()
	if err != nil {
		return err
	}

	// Response
	err = c.JSON(teachers)
	return err
}

func (a *API) UpdateStaffByName(c *fiber.Ctx) error {
	// Parameter check
	name := c.Params("name")
	if len(name) < 1 {
		return errors.New("parameter is not satisfied")
	}

	// Body check
	teacher := model.Staff{}
	if err := c.BodyParser(&teacher); err != nil {
		return err
	}

	// Pass data into repository
	response, err := a.staffRepo.EditStaffByName(name, &teacher)
	if err != nil {
		return err
	}

	// Response
	err = c.JSON(response)
	return err
}

func (a *API) UpdateStaffById(c *fiber.Ctx) error {
	// Parameter check
	id := c.Params("id")
	if len(id) < 1 {
		return errors.New("parameter is not satisfied")
	}

	// Body check
	teacher := model.Staff{}
	if err := c.BodyParser(&teacher); err != nil {
		return err
	}

	// Convert into ObjectID
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Pass data into repository
	response, err := a.staffRepo.EditStaffById(objectId, &teacher)
	if err != nil {
		return err
	}

	// Response
	err = c.JSON(response)
	return err
}

func (a *API) InsertTeachTime(c *fiber.Ctx) error {
	id := c.Params("id")
	if util.IsEmpty(id) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "parameter is not satisfied")
	}

	// Parse id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "malformed id")
	}

	// Parse body
	teachDetails := model.TeachTimeDetail{}
	if c.BodyParser(&teachDetails) != nil {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "body is not satisfied")
	}

	objectId, err = a.staffRepo.AddTeachTime(objectId, &teachDetails)
	return SendSuccessResponse(c, fasthttp.StatusCreated, model.NewResponseID(objectId))
}
