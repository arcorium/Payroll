package controller

import (
	"Penggajian/dbutil"
	"Penggajian/repository"
	"Penggajian/util"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-jwt/jwt/v4"
	"github.com/valyala/fasthttp"
	"time"
)

type API struct {
	app    *fiber.App
	db     *dbutil.Database
	config *dbutil.DBConfig

	// Repository
	userRepo     *repository.UserRepository
	staffRepo    *repository.StaffRepository
	payrollRepo  *repository.PayrollRepository
	positionRepo *repository.PositionRepository
	tokenRepo    *repository.TokenRepository
}

func NewAPI(app_ *fiber.App, config_ *dbutil.DBConfig, db_ *dbutil.Database, userRepo_ *repository.UserRepository,
	staffRepo_ *repository.StaffRepository, payrollRepo_ *repository.PayrollRepository, positionRepo_ *repository.PositionRepository,
	tokenRepo_ *repository.TokenRepository) API {
	return API{app: app_, db: db_, config: config_,
		userRepo:     userRepo_,
		staffRepo:    staffRepo_,
		payrollRepo:  payrollRepo_,
		positionRepo: positionRepo_,
		tokenRepo:    tokenRepo_,
	}
}

func (a *API) HandleAPI() {
	// Middleware
	config := cors.ConfigDefault
	config.AllowCredentials = true
	config.AllowHeaders = "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Authorization"
	a.app.Use(cors.New(config))
	a.app.Use(logger.New(logger.Config{Format: "[${time}] [${ip}]:${port} ${status} - ${method} ${path}\n"}))

	// api/v1
	v1 := a.app.Group("/api/v1")

	// api/v1/user
	userApi := v1.Group("/users")
	// Core
	//userApi.Post("/login", a.Login)
	userApi.Post("/login", a.Login)
	userApi.Post("/req-token", a.RequestToken)

	//authApi := v1.Group("", jwtware.New(jwtware.Config{
	//	SigningKey: []byte(a.config.SecretKey),
	//	//KeyFunc:       a.jwtValidateToken,
	//	ErrorHandler:  a.jwtErrorHandler,
	//	SigningMethod: util.JWT_SIGNING_METHOD},
	//), a.validateAuthorization())
	authApi := v1.Group("") // Dummy

	authUserApi := authApi.Group("/users")
	authUserApi = userApi.Group("")
	authUserApi.Post("/logout", a.Logout)
	// Create
	//authUserApi.Use(a.superOnly())
	authUserApi.Post("/", a.RegisterUser)
	// Delete
	authUserApi.Delete("/id/:id", a.RemoveUserById)
	authUserApi.Delete("/name/:name", a.RemoveUserByName)
	// Get
	authUserApi.Get("/id/:id", a.GetUserById)
	authUserApi.Get("/name/:name", a.GetUserByName)
	authUserApi.Get("/", a.GetUsers)
	// Edit
	authUserApi.Put("/id/:id", a.UpdateUserById)
	authUserApi.Put("/name/:name", a.UpdateUserByName)

	// api/v1/teacher
	staffApi := authApi.Group("/staffs")
	// Create
	staffApi.Post("/", a.RegisterStaff)
	// Get
	staffApi.Get("/id/:id", a.GetStaffById)
	staffApi.Get("/", a.GetStaffs)
	// Edit
	staffApi.Put("/", a.UpdateStaffById)
	// Delete
	staffApi.Delete("/id/:id", a.RemoveStaffById)

	updateStaffApi := staffApi.Group("/id/:id")
	// Teach Time
	updateStaffApi.Post("/teach", a.InsertTeachTime)
	//updateStaffApi.Put("/teach/:uuid", a.UpdateTeachTime)
	updateStaffApi.Delete("/teach/:uuid", a.RemoveTeachTime)
	// Saving
	updateStaffApi.Post("/savings", a.InsertSaving)
	//updateStaffApi.Put("/savings/:uuid", a.UpdateSaving)
	updateStaffApi.Delete("/savings/:uuid", a.RemoveSaving)

	payrollApi := authApi.Group("/payrolls")
	// Get
	payrollApi.Get("/date/:months/:years", a.GetPayrolls)
	payrollApi.Get("/id/:serialNumber", a.GetPayrollBySerialNumber)
	// Edit
	//payrollApi.Put("/:id", a.UpdatePayrollById)
	// Delete
	payrollApi.Delete("/date/:months/:years", a.ClearPayroll)
	// Import
	payrollApi.Post("/imports", a.ImportPayrollFromExcel)
	payrollApi.Post("/copy", a.CopyPayroll)
}

func (a *API) Start(address_ string) error {
	return a.app.Listen(address_)
}

func (a *API) jwtValidateToken(token_ *jwt.Token) (interface{}, error) {
	// Validate algorithm
	if token_.Method.Alg() != util.JWT_SIGNING_METHOD {
		return []byte(a.config.SecretKey), errors.New("wrong algorithm used")
	}

	// Validate expiration time
	claims, ok := token_.Claims.(jwt.MapClaims)
	if !ok {
		return []byte(a.config.SecretKey), errors.New("claims malformed")
	}

	// Check expiration date
	expired := int64(claims["exp"].(float64))

	if expired <= time.Now().Unix() {
		return []byte(a.config.SecretKey), errors.New("token expired")
	}

	// Return signing key
	// TODO: Get signing key from database
	return []byte(a.config.SecretKey), nil
}

func (a *API) jwtErrorHandler(c *fiber.Ctx, err_ error) error {
	return SendErrorResponse(c, fasthttp.StatusUnauthorized, err_.Error())
}
