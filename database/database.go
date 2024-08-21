package database

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"ecom.com/app/models"
	"ecom.com/config"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database instance
type Dbinstance struct {
	Db *gorm.DB

	// db function
	Connect func(schemaName string) *gorm.DB // Connect to database

	// get db from context
	GetDB func(c *fiber.Ctx) *gorm.DB // Get database from context
}

var DB Dbinstance

var schemaDBs = make(map[string]*gorm.DB)

// Connect function
func Connect() {
	// because our config function returns a string, we are parsing our      str to int here
	// schema name
	// DisableForeignKeyConstraintWhenMigrating: false,
	db := createConnection("public")
	DB = Dbinstance{
		Db: db,
		Connect: func(schemaName string) *gorm.DB {
			return getSchemaDB(schemaName)
		},
		GetDB: func(c *fiber.Ctx) *gorm.DB {
			// return c.Locals("DB").(*gorm.DB)
			return getSchemaDB("public")
		},
	}
	// Migrate()
}

func createConnection(schemaName string) *gorm.DB {
	p := config.Config("DB_PORT")

	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		fmt.Println("Error parsing str to int")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=require TimeZone=Asia/Kolkata", config.Config("DB_HOST"), config.Config("DB_USER"), config.Config("DB_PASSWORD"), config.Config("DB_NAME"), port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		// NamingStrategy: schema.NamingStrategy{
		// 	TablePrefix:   schemaName + ".",
		// 	SingularTable: true,
		// 	NoLowerCase:   true,
		// 	NameReplacer:  strings.NewReplacer("CID", "Cid"),
		// },
	})
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}
	log.Println("Connected")
	db.Logger = logger.Default.LogMode(logger.Info)
	schemaDBs[schemaName] = db
	return db
}

// Define a function to get the *gorm.DB object for a given schema name
func getSchemaDB(schemaName string) *gorm.DB {
	// Check if the *gorm.DB object for this schema already exists in the map
	if schemaDB, exists := schemaDBs[schemaName]; exists {
		return schemaDB
	}

	// ctx := DB.Db.Statement.Context

	// Create a new *gorm.DB object for this schema
	// schemaDB := DB.Db.WithContext(ctx).Session(&gorm.Session{NewDB: true})

	// Set the schema for this *gorm.DB object
	// err := schemaDB.Exec(fmt.Sprintf("USE `%s`", schemaName)).Error
	// if err != nil {
	// 	return nil
	// }

	// schemaDB.Config.NamingStrategy = schema.NamingStrategy{
	// 	TablePrefix:   schemaName + ".", // schema name
	// 	SingularTable: true,
	// 	NoLowerCase:   true,
	// 	NameReplacer:  strings.NewReplacer("CID", "Cid"),
	// }

	schemaDB := createConnection(schemaName)

	// Store the *gorm.DB object in the map
	// schemaDBs[schemaName] = schemaDB

	// Return the *gorm.DB object
	return schemaDB
}

func Migrate() {
	// log.Println("running migrations")
	db := DB.Connect("public")
	// base modules
	db.AutoMigrate(&models.Category{})
	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.Tag{})

	log.Println("running migrations")
	// db.AutoMigrate(&model.User{})
	// db.AutoMigrate(&institutes.AttendanceType{}, &institutes.BusRoutes{}, &institutes.Events{}, &institutes.Vehicles{})
	// db.AutoMigrate(&institutes.AttendanceType{}, &institutes.StaffType{}, &institutes.Staffs{})
	// db.AutoMigrate(&institutes.Klass{})
	// db.AutoMigrate(&institutes.Sections{})
	// db.AutoMigrate(&institutes.Staff{})
	// db.AutoMigrate(&institutes.Subjects{})
	// db.AutoMigrate(&institutes.Periods{})
	// db.AutoMigrate(&institutes.Timetable{})
	log.Println("completed migrations")
}
