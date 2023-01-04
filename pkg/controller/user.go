package controller

import (
	"Penggajian/pkg/model"
	"Penggajian/pkg/util"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RegisterUser Password in User still raw and will be hashed here
func (a *API) RegisterUser(c *fiber.Ctx) error {
	user := model.User{}
	// Parse body
	if err := c.BodyParser(&user); err != nil {
		return err
	}

	// Check body field
	if util.IsEmpty(user.Username) || util.IsEmpty(user.Password) {
		return errors.New("body is not satisfied")
	}

	// Pass data into repository
	response, err := a.userRepo.AddUser(&user)
	if err != nil {
		return err
	}

	// Response
	err = c.JSON(response)
	return err
}

func (a *API) RemoveUserByName(c *fiber.Ctx) error {
	// Parameter check
	username := c.Params("username")
	if len(username) < 1 {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "parameter is not satisfied")
	}

	// Pass data into repository
	user, err := a.userRepo.RemoveUserByName(username)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusNoContent, "user not found")
	}

	// Response
	return SendSuccessResponse(c, fasthttp.StatusOK, "user "+user.Username+" deleted!", nil)
}

func (a *API) RemoveUserById(c *fiber.Ctx) error {
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
	user, err := a.userRepo.RemoveUserById(objectId)
	if err != nil {
		return err
	}

	// Response
	err = c.JSON(user)
	return err
}

// GetUserByName get user by the username from uri parameter
func (a *API) GetUserByName(c *fiber.Ctx) error {
	// Parameter check
	username := c.Params("username")
	if len(username) < 1 {
		return errors.New("parameter is not satisfied")
	}

	// Pass data into repository
	user, err := a.userRepo.GetUserByName(username)
	if err != nil {
		return err
	}

	// Response
	err = c.JSON(user)
	return err
}

// GetUserById get user by the objectId from parameter
func (a *API) GetUserById(c *fiber.Ctx) error {
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
	user, err := a.userRepo.GetUserById(objectId)
	if err != nil {
		return err
	}

	// Response
	err = c.JSON(user)
	return err
}

func (a *API) GetUsers(c *fiber.Ctx) error {
	// Get data from repository
	users, err := a.userRepo.GetUsers()
	if err != nil {
		return err
	}

	// Response
	err = c.JSON(users)
	return err
}

func (a *API) UpdateUserByName(c *fiber.Ctx) error {
	// Parameter check
	username := c.Params("username")
	if len(username) < 1 {
		return errors.New("parameter is not satisfied")
	}

	// Body check
	user := model.User{}
	if err := c.BodyParser(&user); err != nil {
		return err
	}

	// Pass data into repository
	response, err := a.userRepo.EditUserByName(username, &user)
	if err != nil {
		return err
	}

	// Response
	err = c.JSON(response)
	return err
}

func (a *API) UpdateUserById(c *fiber.Ctx) error {
	// Parameter check
	id := c.Params("id")
	if len(id) < 1 {
		return errors.New("parameter is not satisfied")
	}

	// Body check
	user := model.User{}
	if err := c.BodyParser(&user); err != nil {
		return err
	}

	// Convert into ObjectID
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Pass data into repository
	response, err := a.userRepo.EditUserById(objectId, &user)
	if err != nil {
		return err
	}

	// Response
	err = c.JSON(response)
	return err
}
