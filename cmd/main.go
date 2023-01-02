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

	// Repository
	userRepo := repository.NewUserRepository(&db, "user")
	
	api := controller.NewAPI(&db, &userRepo)
	api.HandleAPI(app)

	log.Println(app.Listen(":8811"))
}
