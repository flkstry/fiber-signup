package main

import (
	"auth/app/models"
	"auth/database"
	"auth/routes"
	"auth/utils"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type App struct {
	*fiber.App

	DB *database.Database
}

func main() {
	app := App{
		App: fiber.New(fiber.Config{}),
	}

	// argon2 params

	p := &utils.ArgonParams{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}

	// connect to the database
	db, err := database.New(&database.DatabaseConfig{
		Driver:   "postgres",
		Host:     "localhost",
		Username: "postgres",
		Password: "wap12345",
		Port:     5432,
		Database: "auth",
	})

	if err != nil {
		log.Panicln(err)
	}

	// change this with env key
	m, err := utils.NewPasetoMaker("12345678901234567890123456789012")

	if err != nil {
		log.Panicln(err)
	}

	if err != nil {
		log.Println("failed to connect to database:", err.Error())
	} else {
		if db == nil {
			log.Println("failed to connect to database: db variable is nil")
		}

		app.DB = db
		err = app.DB.AutoMigrate(&models.User{})
		if err != nil {
			log.Println("failed to automigrate user model:", err.Error())
			return
		}
	}

	// register routes
	api := app.Group("/api")
	apiv1 := api.Group("/v1")

	routes.UserRoutes(apiv1, app.DB, p, m)

	app.Use(func(ctx *fiber.Ctx) error {
		reqToken := ctx.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		reqToken = splitToken[1]

		_, err := m.VerifyToken(reqToken)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  true,
				"message": err.Error(),
			})
		}

		return ctx.Next()
	})
	routes.BookRoutes(apiv1)

	log.Fatal(app.Listen("127.0.0.1:3000"))
}
