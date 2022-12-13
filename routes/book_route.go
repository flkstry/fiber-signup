package routes

import (
	"auth/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func BookRoutes(api fiber.Router) {
	book := api.Group("book")

	book.Post("/", controllers.GetBook())
}
