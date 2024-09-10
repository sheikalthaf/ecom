package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"ecom.com/app/aws"
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

type DBConfig struct {
	Host     string `json:"DB_HOST"`
	Port     string `json:"DB_PORT"`
	User     string `json:"DB_USER"`
	Password string `json:"DB_PASSWORD"`
	DBName   string `json:"DB_NAME"`
}

var dbConfig *DBConfig

func NewDBConfig(host string, port string, user string, password string, dbName string) {
	dbConfig = &DBConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
	}
}

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
	port, err := strconv.ParseUint(dbConfig.Port, 10, 32)
	if err != nil {
		fmt.Println("Error parsing str to int")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=require TimeZone=Asia/Kolkata", dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.DBName, port)
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
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Cart{})

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

func HandleCredentials(accountType string) {
	switch accountType {
	case "aws":
		sm, err := aws.NewSecretsManager("ap-south-1")
		if err != nil {
			log.Fatal(err.Error())
		}

		secret, err := sm.GetSecret("database")
		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println(secret)
		json.Unmarshal([]byte(secret), &dbConfig)
		NewDBConfig(dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName)
	default:
		var dbConfig DBConfig
		dbConfig.Host = config.Config("DB_HOST")
		dbConfig.Port = config.Config("DB_PORT")
		dbConfig.User = config.Config("DB_USER")
		dbConfig.Password = config.Config("DB_PASSWORD")
		dbConfig.DBName = config.Config("DB_NAME")
		NewDBConfig(dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName)
	}
}
