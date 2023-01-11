package main

import (
	"Penggajian/pkg/controller"
	"Penggajian/pkg/dbutil"
	"Penggajian/pkg/repository"
	"github.com/gofiber/fiber/v2"
	"log"
)
import _ "github.com/joho/godotenv/autoload"

func main() {
	app := fiber.New(fiber.Config{})

	conf, err := dbutil.NewConfig("penggajian", "teacher")
	if err != nil {
		log.Fatalln(err)
	}

	db, err := dbutil.Connect(&conf)
	if err != nil {
		log.Fatalln(err)
	}

	// Set secret key
	//util.JWT_REFRESH_SECRET_KEY = conf.SecretKey

	// Repository
	userRepo := repository.NewUserRepository(&db, "user")
	staffRepo := repository.NewStaffRepository(&db, "teacher")
	positionRepo := repository.NewPositionRepository(&db, "position")
	tokenRepo := repository.NewTokenRepository(&db, "token")

	api := controller.NewAPI(app, &conf, &db, &userRepo, &staffRepo, &positionRepo, &tokenRepo)
	api.HandleAPI()

	log.Println(api.Start(conf.BindAddress))
}
