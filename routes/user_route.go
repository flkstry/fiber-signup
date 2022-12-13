package routes

import (
	"auth/app/controllers"
	"auth/database"
	"auth/utils"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(api fiber.Router, db *database.Database, p *utils.ArgonParams, m *utils.PasetoMaker) {
	signup := api.Group("signup")

	signup.Post("/", controllers.SignUp(db, p))

	signin := api.Group("login")

	signin.Post("/", controllers.SignIn(db, p, m))
}
