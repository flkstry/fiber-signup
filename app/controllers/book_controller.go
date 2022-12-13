package controllers

import "github.com/gofiber/fiber/v2"

func GetBook() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  true,
			"message": "Sudah Login",
		})
	}
}
