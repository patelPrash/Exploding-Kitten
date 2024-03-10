package main

import (
	"fmt"
	"log"
	"os"

	"github.com/anshu7sah/kitten-exploding-backend/middlewares"
	"github.com/anshu7sah/kitten-exploding-backend/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func setupRoutes(app *fiber.App){
	app.Post("/signup",routes.Signup)
	app.Post("/login",routes.Login)
	app.Use("/updatescore",middlewares.JWTMiddleware())
	app.Get("/updatescore",routes.Updatescore)
	app.Get("/getallscores",routes.Getallscores)
	// app.Post("/storegamesession",routes.StoreGamesession)
	// app.Get("/getgamesession",routes.Getgamesession)
}

func main(){
	err:=godotenv.Load()
	if(err!=nil){
		fmt.Println(err);
	}
	app:= fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	setupRoutes(app)

	log.Fatal(app.Listen(os.Getenv("APP_PORT")))

}

