package controller

import (
	"Penggajian/pkg/model"
	"errors"
	"github.com/gofiber/fiber/v2"
	"log"
)

func (a *API) Login(c *fiber.Ctx) error {
	user := model.User{}
	var err error
	if err = c.BodyParser(&user); err != nil {
		log.Println("failed to parse body")
		return err
	}
	if (len(user.TeacherName) < 1 && len(user.Username) < 1) || len(user.Password) < 1 {
		return errors.New("body is not satisfied")
	}
	user, err = a.userRepo.ValidateUser(user.Username, user.Password)

	return err
}

// RegisterUser Password in User still raw and will be hashed here
func (a *API) RegisterUser(c *fiber.Ctx) error {
	user := model.User{}
	if err := c.BodyParser(&user); err != nil {
		log.Println(err)
		return err
	}

	id, err := a.userRepo.AddUser(&user)
	if err != nil {
		return err
	}
	err = c.SendString(id.String())
	return err
}

// GetUserByName get user by the username from uri parameter
func (a *API) GetUserByName(c *fiber.Ctx) error {
	username := c.Params("username")
	if len(username) < 1 {
		return errors.New("parameter is not satisfied")
	}

	user, err := a.userRepo.GetUserByName(username)
	if err != nil {
		return err
	}

	err = c.SendString(user.Id.String())
	return nil
}

// GetUserById get user by the objectId from parameter
func (a *API) GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")
	if len(id) < 1 {
		return errors.New("parameter is not satisfied")
	}

	user, err := a.userRepo.GetUserById(id)
	if err != nil {
		return err
	}

	err = c.SendString(user.Id.String())

	return nil
}

func (a *API) GetUsers(c *fiber.Ctx) error {
	users, err := a.userRepo.GetUsers()
	if err != nil {
		return err
	}

	err = c.SendString(users[0].Username)

	return nil
}
