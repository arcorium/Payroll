package controller

import (
	"Penggajian/pkg/dbutil"
	"Penggajian/pkg/repository"
	"Penggajian/pkg/util"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/valyala/fasthttp"
	"time"
)

type API struct {
	app *fiber.App
	db  *dbutil.Database

	// Repository
	userRepo     *repository.UserRepository
	staffRepo    *repository.StaffRepository
	positionRepo *repository.PositionRepository
	tokenRepo    *repository.TokenRepository
}

func NewAPI(app_ *fiber.App, db_ *dbutil.Database, userRepo_ *repository.UserRepository,
	teacherRepo_ *repository.StaffRepository, positionRepo_ *repository.PositionRepository,
	tokenRepo_ *repository.TokenRepository) API {
	return API{app: app_, db: db_, userRepo: userRepo_, staffRepo: teacherRepo_, positionRepo: positionRepo_, tokenRepo: tokenRepo_}
}

func (a *API) HandleAPI() {

	// Middleware
	a.app.Use(logger.New(logger.Config{}))
	a.app.Use(cors.New(cors.Config{}))

	// api/v1
	v1 := a.app.Group("/api/v1")

	// api/v1/user
	userApi := v1.Group("/users")
	// Core
	userApi.Post("/login", a.Login)
	userApi.Post("/req-token", a.RequestToken)
	userApi.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(util.JWT_REFRESH_SECRET_KEY),
		//KeyFunc:       a.jwtValidateToken,
		ErrorHandler:  a.jwtErrorHandler,
		SigningMethod: util.JWT_SIGNING_METHOD},
	), a.validateAuthorization())
	userApi.Post("/logout", a.Logout)
	// Create
	userApi.Use(a.superOnly())
	userApi.Post("/", a.RegisterUser)
	// Delete
	userApi.Delete("/id/:id", a.RemoveUserById)
	userApi.Delete("/name/:name", a.RemoveUserByName)
	// Get
	userApi.Get("/id/:id", a.GetUserById)
	userApi.Get("/name/:name", a.GetUserByName)
	userApi.Get("/", a.GetUsers)
	// Edit
	userApi.Put("/id/:id", a.UpdateUserById)
	userApi.Put("/name/:name", a.UpdateUserByName)

	// api/v1/teacher
	staffApi := v1.Group("/staffs")
	// Create
	staffApi.Post("/", a.RegisterStaff)
	// Delete
	staffApi.Delete("/id/:id:", a.RemoveStaffById)
	staffApi.Delete("/name/:name:", a.RemoveStaffByName)
	// Get
	staffApi.Get("/id/:id:", a.GetStaffById)
	staffApi.Get("/name/:name:", a.GetStaffByName)
	staffApi.Get("/", a.GetStaffs)
	// Edit
	updateStaffApi := staffApi.Put("/id/:id:", a.UpdateStaffById)
	updateStaffApi.Post("/teach", a.InsertTeachTime)
	staffApi.Put("/name/:name:", a.UpdateStaffByName)

	// api/v1/teacher/pos/
	positionApi := staffApi.Group("/positions")
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

func (a *API) Start(address_ string) error {
	return a.app.Listen(address_)
}

func (a *API) jwtValidateToken(token_ *jwt.Token) (interface{}, error) {
	// Validate algorithm
	if token_.Method.Alg() != util.JWT_SIGNING_METHOD {
		return []byte(util.JWT_ACCESS_SECRET_KEY), errors.New("wrong algorithm used")
	}

	// Validate expiration time
	claims, ok := token_.Claims.(jwt.MapClaims)
	if !ok {
		return []byte(util.JWT_ACCESS_SECRET_KEY), errors.New("claims malformed")
	}

	// Check expiration date
	expired := int64(claims["exp"].(float64))

	// Logging
	fmt.Println("Expired Date: {}", time.Unix(expired, 0))
	fmt.Println("Time: {}", time.Unix(time.Now().Unix(), 0))

	if expired <= time.Now().Unix() {
		return []byte(util.JWT_ACCESS_SECRET_KEY), errors.New("token expired")
	}

	// Return signing key
	// TODO: Get signing key from database
	return []byte(util.JWT_ACCESS_SECRET_KEY), nil
}

func (a *API) jwtErrorHandler(c *fiber.Ctx, err_ error) error {
	return SendErrorResponse(c, fasthttp.StatusUnauthorized, err_.Error())
}
