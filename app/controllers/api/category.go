package api

import (
	"ecom.com/app/models"
	"ecom.com/database"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetCategories function
func GetCategories(c *fiber.Ctx) error {
	db := database.DB.GetDB(c)
	var categories = []models.Category{}
	db.Find(&categories).Where("IsDeleted = ?", false)
	return c.JSON(categories)
}

// GetCategory function
func GetCategory(c *fiber.Ctx) error {
	db := database.DB.GetDB(c)
	var category models.Category
	id := uuid.MustParse(c.Params("id"))
	db.First(&category, id)
	return c.JSON(category)
}

// GetCategoryProducts function
func GetCategoryProducts(c *fiber.Ctx) error {
	db := database.DB.GetDB(c)
	products := []models.Product{}
	categoryID := uuid.MustParse(c.Params("id"))
	db.Model(&models.Product{}).Where(&models.Product{CategoryID: categoryID}).Find(&products)
	tagIds := []uuid.UUID{}
	for _, p := range products {
		for _, t := range p.Tags {
			tagIds = append(tagIds, uuid.MustParse(t))
		}
	}
	tags := GetTagByIds(c, tagIds)
	return c.JSON(mapToProductViewModel(products, tags))
}

// CreateCategory function
func CreateCategory(c *fiber.Ctx) error {
	db := database.DB.GetDB(c)
	category := new(models.Category)
	if err := c.BodyParser(category); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	err := db.Create(&category).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create category", "data": err})
	}
	return c.JSON(category)
}

// UpdateCategory function
func UpdateCategory(c *fiber.Ctx) error {
	db := database.DB.GetDB(c)
	category := new(models.Category)
	if err := c.BodyParser(category); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	err := db.Save(&category).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not update category", "data": err})
	}
	return c.JSON(category)
}

// DeleteCategory function, soft delete
func DeleteCategory(c *fiber.Ctx) error {
	// call update function with deleted as true
	db := database.DB.GetDB(c)
	var category models.Category
	id := uuid.MustParse(c.Params("id"))
	db.First(&category, id)
	category.IsDeleted = true
	err := db.Save(&category).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not delete category", "data": err})
	}
	return c.JSON(category)
}
