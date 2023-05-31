package main

import (
	"fmt"
	"log"

	"Penggajian/controller"
	"Penggajian/dbutil"
	"Penggajian/repository"
	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := fiber.New(fiber.Config{})

	conf, err := dbutil.NewConfig("penggajian", "teacher")
	if err != nil {
		log.Fatalln(err)
	}

	db, err := dbutil.Connect(&conf)
	if err != nil {
		str := fmt.Sprint(" Config: ", conf)
		log.Fatalln(err.Error() + str)
	}
	// Repository
	userRepo := repository.NewUserRepository(&db, "user")
	staffRepo := repository.NewStaffRepository(&db, "staff")
	positionRepo := repository.NewPositionRepository(&db, "position")
	tokenRepo := repository.NewTokenRepository(&db, "token")
	payrollRepo := repository.NewPayrollRepository(&db, "payroll")

	api := controller.NewAPI(app, &conf, &db, &userRepo, &staffRepo, &payrollRepo, &positionRepo, &tokenRepo)
	api.HandleAPI()

	log.Println(api.Start(conf.BindAddress))
}
