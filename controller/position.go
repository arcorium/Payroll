package controller

import (
	"Penggajian/model"
	"errors"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *API) RegisterPosition(c *fiber.Ctx) error {
	position := model.Position{}

	// Parse body
	if err := c.BodyParser(&position); err != nil {
		return err
	}

	// Pass data into repository
	response, err := a.positionRepo.AddPosition(position)
	if err != nil {
		return err
	}

	err = c.JSON(response)
	return err
}

func (a *API) RemovePositionById(c *fiber.Ctx) error {
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
	position, err := a.positionRepo.RemovePositionById(objectId)
	if err != nil {
		return err
	}

	// Response
	err = c.JSON(position)
	return err
}

func (a *API) RemovePositionByName(c *fiber.Ctx) error {
	// Parameter check
	name := c.Params("name")
	if len(name) < 1 {
		return errors.New("parameter is not satisfied")
	}

	// Pass data into repository
	position, err := a.positionRepo.RemovePositionByName(name)
	if err != nil {
		return err
	}

	// Response
	err = c.JSON(position)
	return err
}

func (a *API) GetPositionById(c *fiber.Ctx) error {
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
	position, err := a.positionRepo.GetPositionById(objectId)
	if err != nil {
		return err
	}

	// Response
	err = c.JSON(position)
	return err
}

func (a *API) GetPositionByName(c *fiber.Ctx) error {
	// Parameter check
	name := c.Params("name")
	if len(name) < 1 {
		return errors.New("parameter is not satisfied")
	}

	// Pass data into repository
	position, err := a.positionRepo.GetPositionByName(name)
	if err != nil {
		return err
	}

	// Response
	err = c.JSON(position)
	return err
}

func (a *API) GetPositions(c *fiber.Ctx) error {
	// Get data from repository
	positions, err := a.positionRepo.GetPositions()
	if err != nil {
		return err
	}

	// Response
	err = c.JSON(positions[0])
	return err
}

func (a *API) UpdatePositionById(c *fiber.Ctx) error {
	// Parameter check
	id := c.Params("id")
	if len(id) < 1 {
		return errors.New("parameter is not satisfied")
	}

	// Body check
	position := model.Position{}
	if err := c.BodyParser(&position); err != nil {
		return err
	}

	// Convert into ObjectID
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Pass data into repository
	response, err := a.positionRepo.EditPositionById(objectId, position.Name)
	if err != nil {
		return err
	}

	// Response
	err = c.JSON(response)
	return err
}

func (a *API) UpdatePositionByName(c *fiber.Ctx) error {
	// Parameter check
	name := c.Params("id")
	if len(name) < 1 {
		return errors.New("parameter is not satisfied")
	}

	// Body check
	position := model.Position{}
	if err := c.BodyParser(&position); err != nil {
		return err
	}

	// Pass data into repository
	response, err := a.positionRepo.EditPositionByName(name, position.Name)
	if err != nil {
		return err
	}

	// Response
	err = c.JSON(response)
	return err
}
