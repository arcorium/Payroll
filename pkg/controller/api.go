package controller

import (
	"Penggajian/pkg/dbutil"
	"Penggajian/pkg/repository"
	"Penggajian/pkg/util"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v3"
)

type API struct {
	db *dbutil.Database

	// Repository
	userRepo     *repository.UserRepository
	staffRepo    *repository.StaffRepository
	positionRepo *repository.PositionRepository
	teachRepo    *repository.TeachRepository
	tokenRepo    *repository.TokenRepository
}

func NewAPI(db_ *dbutil.Database, userRepo_ *repository.UserRepository,
	teacherRepo_ *repository.StaffRepository, positionRepo_ *repository.PositionRepository,
	teachRepo_ *repository.TeachRepository, tokenRepo_ *repository.TokenRepository) API {
	return API{db: db_, userRepo: userRepo_, staffRepo: teacherRepo_, positionRepo: positionRepo_,
		teachRepo: teachRepo_, tokenRepo: tokenRepo_}
}

func (a *API) HandleAPI(app_ *fiber.App) {

	// Middleware
	app_.Use(logger.New(logger.Config{}))

	// api/v1
	api := app_.Group("/api/v1")
	api.Post("/login", a.Login)
	api.Post("/req_token", a.RequestToken)

	api.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(util.JWT_REFRESH_SECRET_KEY),
		SigningMethod: util.JWT_SIGNING_METHOD},
	))

	api.Post("/logout", a.Logout)
	// api/v1/user
	userApi := api.Group("/user")
	// Create
	userApi.Post("/", a.RegisterUser)
	// Delete
	userApi.Delete("/id/:id:", a.RemoveUserById)
	userApi.Delete("/name/:username:", a.RemoveUserByName)
	// Get
	userApi.Get("/id/:id:", a.GetUserById)
	userApi.Get("/name/:username:", a.GetUserByName)
	userApi.Get("/", a.GetUsers)
	// Edit
	userApi.Put("/id/:id:", a.UpdateUserById)
	userApi.Put("/name/:username:", a.UpdateUserByName)

	// api/v1/teacher
	teacherApi := api.Group("/teacher")
	// Create
	teacherApi.Post("/", a.RegisterStaff)
	// Delete
	teacherApi.Delete("/id/:id:", a.RemoveStaffById)
	teacherApi.Delete("/name/:username:", a.RemoveStaffByName)
	// Get
	teacherApi.Get("/id/:id:", a.GetStaffById)
	teacherApi.Get("/name/:name:", a.GetStaffByName)
	teacherApi.Get("/", a.GetStaffs)
	// Edit
	teacherApi.Put("/id/:id:", a.UpdateStaffById)
	teacherApi.Put("/name/:name:", a.UpdateStaffByName)

	// api/v1/teacher/pos/
	positionApi := teacherApi.Group("/pos")
	// Create
	positionApi.Post("/", a.RegisterPosition)
	// Delete
	positionApi.Delete("/id/:id:", a.RemovePositionById)
	positionApi.Delete("/name/:name:", a.RemovePositionByName)
	// Get
	positionApi.Get("/id/:id:", a.GetPositionById)
	positionApi.Get("/name/:name:", a.GetPositionByName)
	positionApi.Get("/", a.GetPositions)
	// Edit
	positionApi.Put("/id/:id:", a.UpdatePositionById)
	positionApi.Put("/name/:name:", a.UpdatePositionByName)
}
