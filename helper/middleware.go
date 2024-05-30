package helpers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	models "github.com/NdnHnnt/projectSprint_BeliMang/model"

	// "HaloSuster/db"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gofiber/fiber/v2"
)

type UserClaims struct {
	jwt.RegisteredClaims
	Id    string `json:"id"`
	Email string `json:"email"`
}

func SignUserJWT(user models.UserModel) string {
	// expiredIn := 28800 // 8 hours
	exp := time.Now().Add(time.Hour * 8)
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(exp),
			Issuer:    "BeliMang",
		},
		Id:    user.ID,
		Email: user.Email,
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)
	fmt.Println("Email:", user.Email)
	jwtSecret := os.Getenv("JWT_SECRET")
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return ""
	}
	return signedToken
}

func ParseToken(jwtToken string) (string, string, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, OK := token.Method.(*jwt.SigningMethodHMAC); !OK {
			return nil, errors.New("bad signed method received")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return "", "", err
	}
	parsedToken, OK := token.Claims.(jwt.MapClaims)
	if !OK {
		return "", "", errors.New("unable to parse claims")
	}
	id := fmt.Sprint(parsedToken["id"])       // changed "Id" to "id"
	email := fmt.Sprint(parsedToken["email"]) // changed "Email" to "email"
	return id, email, nil
}

func getBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}

func AuthAdminMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(http.StatusUnauthorized).SendString("Missing Authorization header")
	}

	// Extract the JWT token from the Authorization header
	tokenStr, err := getBearerToken(authHeader)

	if err != nil {
		return c.Status(http.StatusUnauthorized).SendString("Invalid Authorization header format")
	}

	// Parse and validate the JWT token, and extract the Nip
	id, email, err := ParseToken(tokenStr)
	if err != nil {
		fmt.Println("Error parsing token: ", err)
		return c.Status(http.StatusUnauthorized).SendString("Invalid JWT token")
	}

	// Check if email and id is in admin table
	exists, err := ValidateAdmin(email, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}
	if !exists {
		return c.Status(http.StatusUnauthorized).SendString("Unauthorized")
	}

	// Store the Nip in the request context
	c.Locals("userEmail", email)
	c.Locals("userId", id)

	// Continue with the next middleware function or the request handler
	return c.Next()
}

func AuthUserMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(http.StatusUnauthorized).SendString("Missing Authorization header")
	}

	// Extract the JWT token from the Authorization header
	tokenStr, err := getBearerToken(authHeader)
	fmt.Println("Token:", tokenStr)
	if err != nil {
		return c.Status(http.StatusUnauthorized).SendString("Invalid Authorization header format")
	}

	// Parse and validate the JWT token, and extract the Nip
	id, email, err := ParseToken(tokenStr)
	fmt.Println("Email:", email)
	if err != nil {
		return c.Status(http.StatusUnauthorized).SendString("Invalid JWT token")
	}

	// Check if email and id is in admin table
	exists, err := ValidateUser(email, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	if !exists {
		return c.Status(http.StatusUnauthorized).SendString("Unauthorized")
	}

	// Store the Nip in the request context
	c.Locals("userEmail", email)
	c.Locals("userId", id)

	// Continue with the next middleware function or the request handler
	return c.Next()
}

func CheckPassword(password, dbpassword string) bool {
	return password == dbpassword
}


