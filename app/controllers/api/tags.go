package api

import (
	"ecom.com/app/models"
	"ecom.com/database"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetTags function
func GetTags(c *fiber.Ctx) error {
	db := database.DB.GetDB(c)
	var tags = []models.Tag{}
	db.Find(&tags).Where("IsDeleted = ?", false)
	return c.JSON(tags)
}

// GetTag function
func GetTag(c *fiber.Ctx) error {
	db := database.DB.GetDB(c)
	var tag models.Tag
	id := uuid.MustParse(c.Params("id"))
	db.First(&tag, id)
	return c.JSON(tag)
}

func GetTagByIds(c *fiber.Ctx, ids []uuid.UUID) []models.Tag {
	if len(ids) == 0 {
		return []models.Tag{}
	}
	db := database.DB.GetDB(c)
	var tags = []models.Tag{}
	db.Find(&tags, ids)
	return tags
}

// CreateTag function
func CreateTag(c *fiber.Ctx) error {
	db := database.DB.GetDB(c)
	tag := new(models.Tag)
	if err := c.BodyParser(tag); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	err := db.Create(&tag).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create tag", "data": err})
	}
	return c.JSON(tag)
}

// UpdateTag function
func UpdateTag(c *fiber.Ctx) error {
	db := database.DB.GetDB(c)
	tag := new(models.Tag)
	if err := c.BodyParser(tag); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	err := db.Save(&tag).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not update tag", "data": err})
	}
	return c.JSON(tag)
}

// DeleteTag function
func DeleteTag(c *fiber.Ctx) error {
	db := database.DB.GetDB(c)
	var tag models.Tag
	id := uuid.MustParse(c.Params("id"))
	db.First(&tag, id)
	tag.IsDeleted = true
	err := db.Save(&tag).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not delete tag", "data": err})
	}
	return c.JSON(tag)
}
