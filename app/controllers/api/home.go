package api

import (
	"ecom.com/app/models"
	"ecom.com/database"
	"github.com/gofiber/fiber/v2"
)

// GetCategories and 10 top products of each category
func GetHome(c *fiber.Ctx) error {
	db := database.DB.GetDB(c)
	var categories = []models.Category{}
	db.Find(&categories).Where("IsDeleted = ?", false)

	for i := range categories {
		products := []models.Product{}
		db.Model(models.Product{}).Limit(6).Find(&products, `"CategoryId" = ?`, categories[i].ID)
		categories[i].Product = products
	}
	return c.JSON(categories)
}
