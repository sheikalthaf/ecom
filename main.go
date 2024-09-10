package main

import (
	"fmt"

	"ecom.com/app/aws"
	"ecom.com/app/local"
	"ecom.com/app/storage"
	"ecom.com/config"
	"ecom.com/database"
	"ecom.com/router"
	"ecom.com/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// load the env variables
	config.LoadEnv()

	// create a new fiber app
	app := fiber.New()

	// create a new fiber app with tcp6 network
	// app := fiber.New(fiber.Config{
	// 	Network: fiber.NetworkTCP6,
	// })

	// Default middleware config
	app.Use(compress.New())

	// cors middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	accountType := "aws"

	// secret manager
	database.HandleCredentials(accountType)
	aws.InitS3Config()
	var storage storage.Storage

	switch accountType {
	case "aws":

		// Initialize the storage handler
		config, err := aws.LoadConfig()
		if err != nil {
			// log.Fatalf("Failed to load AWS config: %v", err)
		}
		storage = aws.NewS3Storage(config)
	default:
		storage = local.NewLocalImageStorage()
	}

	// initialize the storage handler
	utilities.NewHandler(storage)

	// connect to the database
	database.Connect()

	// utilities.ImageInit()
	// err := utilities.InitS3("ssr-bkt")
	// if err != nil {
	// 	// log.Fatalf("Failed to initialize S3: %v", err)
	// }

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

	port := config.Config("PORT")
	if port == "" {
		port = "5000" // fallback port
	}
	fmt.Println("port", port)
	app.Listen(":" + port)
}
