package controllers

import (
	"auth/app/models"
	"auth/database"
	"auth/utils"
	"errors"
	"fmt"
	"log"
	"net/mail"
	"time"

	"github.com/go-passwd/validator"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SignUp(db *database.Database, p *utils.ArgonParams) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// create a user
		User := new(models.User)

		// parse body to get "User" data
		if err := ctx.BodyParser(User); err != nil {
			log.Println("An error occurred when parsing the new user: " + err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}

		// Email Validation

		// check if user already type email
		if len(User.Email) == 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  false,
				"message": "Email belum dimasukan",
			})
		}

		// check if email is in approriate format
		_, err := mail.ParseAddress(User.Email)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  false,
				"message": "Format email yang anda masukan salah",
			})
		}

		// check if email is already registered
		var CheckUser models.User
		res := db.Where("email = ?", User.Email).Find(&CheckUser)
		if res.RowsAffected > 0 {
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"status":  false,
				"message": fmt.Sprintf("email %s sudah terdaftar", User.Email),
			})
		}

		// Name validation

		// check if user already typing name
		if len(User.Name) == 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  false,
				"message": "Nama belum dimasukan",
			})
		}

		// validate password requirement
		pwdValidator := validator.New(validator.MinLength(8, errors.New("password minimal 8 digits")))
		err = pwdValidator.Validate(User.Password)
		if err != nil {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  false,
				"message": "password minimal 8 digits",
			})
		}

		// generate hashed password
		hash, err := utils.GenerateHashedPassword(User.Password, p)
		if err != nil {
			log.Println("File to generate hashed password: " + err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}

		// save data to DB
		User.Password = string(hash[:])
		if res := db.Create(&User); res.Error != nil {
			log.Println("An error occurred when storing the new user: " + res.Error.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  false,
				"message": res.Error.Error(),
			})
		}

		// return response
		return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status":  true,
			"message": "Account created",
		})
	}
}

func SignIn(db *database.Database, p *utils.ArgonParams, m *utils.PasetoMaker) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// create a user
		User := new(models.User)

		// parse body to get "User" data
		if err := ctx.BodyParser(User); err != nil {
			log.Println("An error occurred when parsing the new user: " + err.Error())
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		// get hashed password from database
		// check user by email
		SavedUser := new(models.User)
		res := db.Where("email = ?", User.Email).Take(&SavedUser).Error
		if errors.Is(res, gorm.ErrRecordNotFound) {
			// return response
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  false,
				"message": res.Error(),
			})
		}

		// generate hashed password
		match, err := utils.ComparePasswordAndHash(User.Password, SavedUser.Password)
		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		if match {
			accessToken, err := m.CreateToken(User.Email, time.Duration(time.Minute*1))
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  false,
					"message": err.Error(),
				})
			}
			// return response
			return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
				"status":       true,
				"message":      "Berhasil login",
				"access_token": accessToken,
			})
		}

		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Password yang anda masukan salah",
		})
	}
}
