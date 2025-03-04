package main

import (
	"ecom.com/database"
	"ecom.com/router"
	"ecom.com/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	database.Connect()
	app := fiber.New(fiber.Config{
		Network: fiber.NetworkTCP6,
	})

	// Default middleware config
	app.Use(compress.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// utilities.ImageInit()
	err := utilities.InitS3("ssr-bkt")
	if err != nil {
		// log.Fatalf("Failed to initialize S3: %v", err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Static("/api/images", "./images", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		MaxAge:        3600,
		CacheDuration: 3600,
	})
	app.Static("/api/audio", "./audio")
	router.PublicRoutes(app)

	// JWT Middleware
	// app.Use(jwt.New(jwt.Config{
	// 	SigningKey: jwt.SigningKey{Key: utilities.AppConfig.JWTSecret},
	// }))

	// app.Use(func(c *fiber.Ctx) error {
	// 	// Get the schema name from the request parameters
	// 	jwt := handler.TokenValues(c)

	// 	if jwt == nil {
	// 		return c.Next()
	// 	}

	// 	// schemaName := jwt.DbName

	// 	// // Get the *gorm.DB object for this schema
	// 	// schemaDB := database.DB.Connect(schemaName)
	// 	// if err != nil {
	// 	//     return c.Status(http.StatusInternalServerError).SendString(err.Error())
	// 	// }

	// 	// Store the *gorm.DB object in the context object
	// 	// c.Context().SetUserValue("schemaDB", schemaDB)
	// 	return c.Next()
	// })

	router.PrivateRoutes(app)
	// handle unavailable route
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

	app.Listen(":8101")
}
