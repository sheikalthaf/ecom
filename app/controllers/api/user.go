package api

import (
	"ecom.com/app/models"
	"ecom.com/database"
	"github.com/gofiber/fiber/v2"

	"github.com/google/uuid"
)

// Register user
func Register(c *fiber.Ctx) error {
	db := database.DB.GetDB(c)
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return err
	}
	db.Create(&user)
	return c.JSON(user)
}

// GetUser function
func GetUser(c *fiber.Ctx) error {
	db := database.DB.GetDB(c)
	var user UserViewModel
	id := uuid.MustParse(c.Params("id"))
	db.First(&user, id)
	return c.JSON(user)
}

// update user without password
func UpdateUser(c *fiber.Ctx) error {
	db := database.DB.GetDB(c)
	user := new(UserViewModel)
	if err := c.BodyParser(user); err != nil {
		return err
	}

	id := uuid.MustParse(c.Params("id"))
	db.Model(&UserViewModel{}).Where("id = ?", id).Updates(&user)
	return c.JSON(user)
}

// get all users
func GetAllUsers(c *fiber.Ctx) error {
	db := database.DB.GetDB(c)
	var users []UserViewModel
	db.Find(&users)
	return c.JSON(users)
}

type UserViewModel struct {
	ID       string `json:"Id"`
	Name     string `json:"Name"`
	Email    string `json:"Email"`
	MobileNo string `json:"MobileNo"`
}
