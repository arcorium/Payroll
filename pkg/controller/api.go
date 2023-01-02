package controller

import (
	"Penggajian/pkg/dbutil"
	"Penggajian/pkg/repository"
	"github.com/gofiber/fiber/v2"
)

type API struct {
	db *dbutil.Database

	// Repository
	userRepo *repository.UserRepository
}

func NewAPI(db_ *dbutil.Database, userRepo_ *repository.UserRepository) API {
	return API{db: db_, userRepo: userRepo_}
}

func (a *API) HandleAPI(app_ *fiber.App) {

	// api/v1
	api := app_.Group("/api/v1")
	api.Post("/login", a.Login)

	// api/v1/user
	userApi := api.Group("/user")
	userApi.Get("/", a.GetUsers)
	userApi.Get("/name/:username:", a.GetUserByName)
	userApi.Get("/id/:id:", a.GetUserById)
	userApi.Post("/", a.RegisterUser)

	//app_.Use(jwtware.New(jwtware.Config{
	//	SigningKey: []byte("secret"),
	//}))
}
