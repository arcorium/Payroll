package controller

import (
	"Penggajian/model"
	"Penggajian/util"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RegisterUser Password in User still raw and will be hashed here
func (a *API) RegisterUser(c *fiber.Ctx) error {
	user := model.User{}
	// Parse body
	if err := c.BodyParser(&user); err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, "failed to parse body")
	}

	// Check body field
	if util.IsEmpty(user.Username) || util.IsEmpty(user.Password) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_BODY)
	}

	// Pass data into repository
	id, err := a.userRepo.AddUser(&user)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "user already registered")
	}

	// Get inserted user
	user, err = a.userRepo.GetUserById(id)
	if util.IsError(err) {
		return SendErrorResponse(c, fasthttp.StatusNotFound, "there is no user found")
	}

	// Response
	return SendSuccessResponse(c, fiber.StatusCreated, user)
}

func (a *API) RemoveUserByName(c *fiber.Ctx) error {
	// Parameter check
	username := c.Params("name")
	if len(username) < 1 {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "parameter is not satisfied")
	}

	// Pass data into repository
	user, err := a.userRepo.RemoveUserByName(username)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusNotFound, "user not found")
	}

	// Response
	return SendSuccessResponse(c, fasthttp.StatusOK, model.NewResponseID(user.Id))
}

func (a *API) RemoveUserById(c *fiber.Ctx) error {
	// Parameter check
	id := c.Params("id")
	if len(id) < 1 {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "parameter is not satisfied")
	}

	// Convert into ObjectID
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "malformed id")
	}

	// Pass data into repository
	user, err := a.userRepo.RemoveUserById(objectId)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusNotModified, "user not removed")
	}

	// Response
	return SendSuccessResponse(c, fasthttp.StatusOK, model.NewResponseID(user.Id))
}

// GetUserByName get user by the username from uri parameter
func (a *API) GetUserByName(c *fiber.Ctx) error {
	// Parameter check
	username := c.Params("name")
	if len(username) < 1 {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "parameter is not satisfied")
	}

	// Pass data into repository
	user, err := a.userRepo.GetUserByName(username)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusNotFound, "user not found")
	}

	// Response
	return SendSuccessResponse(c, fasthttp.StatusOK, user)
}

// GetUserById get user by the objectId from parameter
func (a *API) GetUserById(c *fiber.Ctx) error {
	// Parameter check
	id := c.Params("id")
	if len(id) < 1 {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "parameter is not satisfied")
	}

	// Convert into ObjectID
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "malformed id")
	}

	// Pass data into repository
	user, err := a.userRepo.GetUserById(objectId)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusNotFound, "user not found")
	}

	// Response
	return SendSuccessResponse(c, fasthttp.StatusOK, user)
}

func (a *API) GetUsers(c *fiber.Ctx) error {
	// Get data from repository
	users, err := a.userRepo.GetUsers()
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, "can't get users")
	}

	// Response
	return SendSuccessResponse(c, fasthttp.StatusOK, users)
}

func (a *API) UpdateUserByName(c *fiber.Ctx) error {
	// Parameter check
	username := c.Params("name")
	if util.IsEmpty(username) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "parameter is not satisfied")
	}

	// Body check
	user := model.User{}
	if err := c.BodyParser(&user); err != nil {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "body is not satisfied")
	}

	// Pass data into repository
	err := a.userRepo.EditUserByName(username, &user)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusNotModified, "user not found")
	}

	// Response
	return SendSuccessResponse(c, fasthttp.StatusOK, nil)
}

func (a *API) UpdateUserById(c *fiber.Ctx) error {
	// Parameter check
	id := c.Params("id")
	if util.IsEmpty(id) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "parameter is not satisfied")
	}

	// Body check
	user := model.User{}
	if err := c.BodyParser(&user); err != nil {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "body is not satisfied")
	}

	// Convert into ObjectID
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "malformed id")
	}

	// Pass data into repository
	err = a.userRepo.EditUserById(objectId, &user)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusNotModified, "user not found")
	}

	// Get inserted user
	user, err = a.userRepo.GetUserById(objectId)
	if util.IsError(err) {
		return SendErrorResponse(c, fasthttp.StatusNotFound, "there is no user found")
	}

	// Response
	return SendSuccessResponse(c, fasthttp.StatusOK, user)
}
