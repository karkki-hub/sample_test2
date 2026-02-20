package main

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Login route to generate JWT
	e.POST("/login", func(c echo.Context) error {
		type Login struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		var login Login
		if err := c.Bind(&login); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		}

		// Validate credentials (simplified for demo)
		if login.Username != "user" || login.Password != "pass" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		}

		// Create JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": login.Username,
			"exp":      time.Now().Add(time.Hour * 1).Unix(),
		})
		tokenString, err := token.SignedString([]byte("my-secret"))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
		}

		return c.JSON(http.StatusOK, map[string]string{"token": tokenString})
	})

	// Protected routes with JWT middleware
	protected := e.Group("/secure", middleware.JWT([]byte("my-secret")))
	protected.GET("/profile", func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		username := claims["username"].(string)
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello, " + username})
	})

	e.Logger.Fatal(e.Start(":8080"))
}
