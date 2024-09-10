package api

import (
	"strconv"

	"ecom.com/app/models"
	"ecom.com/database"
	"ecom.com/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var imageFolder = "product"

// GetProducts function
func GetProducts(c *fiber.Ctx) error {
	db := database.DB.GetDB(c)
	var products = []models.Product{}
	db.Model(&models.Product{}).Joins("Category").Find(&products)

	tagIds := []uuid.UUID{}
	for _, p := range products {
		for _, t := range p.Tags {
			tagIds = append(tagIds, uuid.MustParse(t))
		}
	}
	tags := GetTagByIds(c, tagIds)
	return c.JSON(mapToProductViewModel(products, tags))
}

// get products by category id
func GetProductsByCategoryID(c *fiber.Ctx) error {
	db := database.DB.GetDB(c)
	var products = []models.Product{}
	categoryID := uuid.MustParse(c.Params("id"))
	db.Model(&models.Product{}).Where("CategoryID = ?", categoryID).Find(&products)

	tagIds := []uuid.UUID{}
	for _, p := range products {
		for _, t := range p.Tags {
			tagIds = append(tagIds, uuid.MustParse(t))
		}
	}
	tags := GetTagByIds(c, tagIds)
	return c.JSON(mapToProductViewModel(products, tags))
}

// GetProduct function
func GetProduct(c *fiber.Ctx) error {
	db := database.DB.GetDB(c)
	var product models.Product
	id := uuid.MustParse(c.Params("id"))
	db.First(&product, id)
	return c.JSON(product)
}

// CreateProduct function
func CreateProduct(c *fiber.Ctx) error {
	db := database.DB.GetDB(c)
	product := new(models.Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	newImages := []string{}
	oldImages := []string{}

	for i, p := range product.Image {
		oldImage, newImage, err := utilities.Image.UploadImage(c, "Image"+strconv.Itoa(i), imageFolder, p)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
		}
		newImages = append(newImages, newImage)
		oldImages = append(oldImages, *oldImage)
	}
	product.Image = newImages
	err := db.Create(&product).Error
	// delete the old images
	for _, oldImage := range oldImages {
		utilities.Image.DeleteImage(oldImage, &oldImage, imageFolder, err != nil)
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create product", "data": err})
	}
	return c.JSON(product)
}

// UpdateProduct function
func UpdateProduct(c *fiber.Ctx) error {
	db := database.DB.GetDB(c)
	product := new(models.Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	var oldProduct models.Product
	id := uuid.MustParse(c.Params("id"))
	db.First(&oldProduct, id)
	if oldProduct.Image[0] != product.Image[0] {
		oldImage := oldProduct.Image[0]
		oldProduct.Image = product.Image
		// delete the old images
		utilities.Image.DeleteImage(product.Image[0], &oldImage, imageFolder, false)
	}
	if err := db.Save(&product).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not update user", "data": err})
	}
	return c.JSON(product)
}

// DeleteProduct function
func DeleteProduct(c *fiber.Ctx) error {
	db := database.DB.GetDB(c)
	var product models.Product
	id := uuid.MustParse(c.Params("id"))
	db.First(&product, id)
	if err := db.Delete(&product, product.ID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not delete product", "data": err})
	}
	// delete the old images
	for _, oldImage := range product.Image {
		utilities.Image.DeleteImage(oldImage, &oldImage, imageFolder, false)
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User successfully deleted"})
}

// map to view model
func mapToProductViewModel(product []models.Product, tags []models.Tag) []ProductViewModel {
	tagsMap := map[uuid.UUID]string{}
	for _, t := range tags {
		tagsMap[t.ID] = t.Name
	}
	var productViewModel = []ProductViewModel{}
	for _, p := range product {
		if p.Tags == nil {
			p.Tags = []string{}
		}
		var tagNames = []string{}
		for _, t := range p.Tags {
			tagNames = append(tagNames, tagsMap[uuid.MustParse(t)])
		}
		productViewModel = append(productViewModel, ProductViewModel{
			ID:           p.ID,
			Name:         p.Name,
			CategoryID:   p.CategoryID,
			CategoryName: p.Category.Name,
			ProductCode:  p.ProductCode,
			Price:        p.Price,
			Description:  p.Description,
			Size:         p.Size,
			Quantity:     p.Quantity,
			IsDeleted:    p.IsDeleted,
			Image:        p.Image,
			Tags:         p.Tags,
			TagNames:     tagNames,
		})
	}
	return productViewModel
}

// ProductViewModel struct
type ProductViewModel struct {
	ID           uuid.UUID `json:"Id"`
	Name         string    `json:"Name"`
	CategoryID   uuid.UUID `json:"CategoryId"`
	CategoryName string    `json:"CategoryName"`
	ProductCode  string    `json:"ProductCode"`
	Price        float64   `json:"Price"`
	Description  string    `json:"Description"`
	Size         string    `json:"Size"`
	Quantity     int       `json:"Quantity"`
	IsDeleted    bool      `json:"IsDeleted"`
	Image        []string  `json:"Image"`
	Tags         []string  `json:"Tags"`
	TagNames     []string  `json:"TagNames"`
}
