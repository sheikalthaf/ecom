package utilities

import (
	"strconv"

	"ecom.com/config"
)

type Config struct {
	// DB
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	// JWT
	JWTSecret []byte
	// Server
	ServerPort string
	// Email
	EmailHost     string
	EmailPort     int
	EmailUsername string
	EmailPassword string
	EmailFrom     string
	// Firebase
	FireServerKey string
}

var AppConfig Config

func InitConfig() {
	emailPort, _ := strconv.Atoi(config.Config("EMAIL_PORT"))
	jwtSecret := []byte(config.Config("JWT_SECRET_KEY"))
	AppConfig = Config{
		// DB
		DBHost:     config.Config("DB_HOST"),
		DBPort:     config.Config("DB_PORT"),
		DBUser:     config.Config("DB_USER"),
		DBPassword: config.Config("DB_PASSWORD"),
		DBName:     config.Config("DB_NAME"),
		// JWT
		JWTSecret: jwtSecret,
		// Email
		EmailHost:     config.Config("EMAIL_HOST"),
		EmailPort:     emailPort,
		EmailUsername: config.Config("EMAIL_USERNAME"),
		EmailPassword: config.Config("EMAIL_PASSWORD"),
		EmailFrom:     config.Config("EMAIL_FROM"),
		// Firebase
		FireServerKey: config.Config("FIREBASE_SERVER_KEY"),
	}
}
