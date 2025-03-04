package router

import (
	"ecom.com/app/controllers/api"
	"github.com/gofiber/fiber/v2"
)

func PublicRoutes(app *fiber.App) {
	// home
	app.Get("/home", api.GetHome)

	// get products
	app.Get("/products", api.GetProducts)
	app.Get("/products/:id", api.GetProduct)

	// get categories
	app.Get("/categories", api.GetCategories)
	app.Get("/categories/:id", api.GetCategory)
	app.Get("/categories/products/:id", api.GetCategoryProducts)
}

func PrivateRoutes(app *fiber.App) {
	// products
	app.Post("/products", api.CreateProduct)
	app.Put("/products/:id", api.UpdateProduct)
	app.Delete("/products/:id", api.DeleteProduct)

	// categories
	app.Post("/categories", api.CreateCategory)
	app.Put("/categories/:id", api.UpdateCategory)
	app.Delete("/categories/:id", api.DeleteCategory)

	// tags
	app.Get("/tags", api.GetTags)
	app.Get("/tags/:id", api.GetTag)
	app.Post("/tags", api.CreateTag)
	app.Put("/tags/:id", api.UpdateTag)
	app.Delete("/tags/:id", api.DeleteTag)
}
