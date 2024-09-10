package api

import (
	"ecom.com/app/models"
	"ecom.com/database"
	"github.com/gofiber/fiber/v2"
)

// Get cart by user id
func GetCartByUserID(c *fiber.Ctx) error {
	db := database.DB.GetDB(c)
	userID := c.Params("id")
	var cart []CartViewModel

	db.Table("carts").
		Select("carts.id, carts.user_id, carts.product_id, products.name as product_name, categories.name as category_name, products.price, products.image as product_image, carts.quantity, carts.created_at").
		Joins("JOIN products ON products.id = carts.product_id").
		Joins("JOIN categories ON categories.id = products.category_id").
		Where("carts.user_id = ?", userID).
		Scan(&cart)

	return c.JSON(cart)
}

// Add product to cart
func AddProductToCart(c *fiber.Ctx) error {
	db := database.DB.GetDB(c)
	cart := new(models.Cart)
	if err := c.BodyParser(cart); err != nil {
		return err
	}
	db.Create(&cart)
	return c.JSON(cart)
}

// Update product in cart
func UpdateProductInCart(c *fiber.Ctx) error {
	db := database.DB.GetDB(c)
	cart := new(models.Cart)
	if err := c.BodyParser(cart); err != nil {
		return err
	}
	db.Save(&cart)
	return c.JSON(cart)
}

// Delete product from cart
func DeleteProductFromCart(c *fiber.Ctx) error {
	db := database.DB.GetDB(c)
	id := c.Params("id")
	db.Delete(&models.Cart{}, id)
	return nil
}

type CartViewModel struct {
	ID           string `json:"Id"`
	UserID       string `json:"UserId"`
	ProductID    string `json:"ProductId"`
	ProductName  string `json:"ProductName"`
	CategoryName string `json:"CategoryName"`
	Price        int    `json:"Price"`
	ProductImage string `json:"ProductImage"`
	Quantity     int    `json:"Quantity"`
	CreatedAt    string `json:"CreatedAt"`
}
