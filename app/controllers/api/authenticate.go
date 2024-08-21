package api

import (
	"encoding/json"
	"errors"
	"time"

	"ecom.com/app/models"
	"ecom.com/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func AuthenticateUser(c *fiber.Ctx) error {
	login := new(Login)
	err := c.BodyParser(login)

	// throw error for invalid json
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Invalid JSON", "data": err})
	}

	// check whether the user is valid or not
	user, error := getUser(login.UserName, login.Password)

	// Throws Unauthorized error
	if error {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	t, err := generateAccessToken(*user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not generate token", "data": err})
	}

	return c.JSON(fiber.Map{"token": LoginResponse{AccessToken: *t}})
}

func getUser(userName string, password string) (*models.UserDetails, bool) {
	if userName == "john" && password == "doe" {
		details := models.UserDetails{
			ID:   uuid.New(),
			Name: "John Doe",
		}
		return &details, false
	} else {
		return nil, true
	}
}

func generateAccessToken(user models.UserDetails) (*string, error) {
	// var year = "2022"
	userDetails := new(UserTokenDetails)
	userDetails.UserId = user.ID
	userDetails.Name = user.Name

	t, err := GenerateJWT(userDetails, time.Hour*30)
	if err != nil {
		return nil, errors.New("could not generate token")
	}

	return &t, nil
}

func GenerateJWT(details interface{}, duration time.Duration) (string, error) {
	// Create the Claims
	claims := jwt.MapClaims{
		"name":    "John Doe",
		"admin":   true,
		"exp":     time.Now().Add(duration).Unix(),
		"details": details,
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString(utilities.AppConfig.JWTSecret)
	if err != nil {
		return "", err
	}
	return t, nil
}

func TokenValues(c *fiber.Ctx) *UserTokenDetails {
	user, err := c.Locals("user").(*jwt.Token)
	if !err {
		return nil
	}
	claims := user.Claims.(jwt.MapClaims)

	s := UserTokenDetails{}
	jsonString, _ := json.Marshal(claims["details"])
	json.Unmarshal(jsonString, &s)
	return &s
}

type LoginResponse struct {
	AccessToken string
}

type Login struct {
	UserName string
	Password string
}

type UserTokenDetails struct {
	UserId    uuid.UUID
	RoleId    uuid.UUID
	UserRefId uuid.UUID
	Name      string
}
