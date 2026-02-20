package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	//"github.com/labstack/echo/v4/middleware"
	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
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
		fmt.Println(token.Raw)

		return c.JSON(http.StatusOK, map[string]string{"token": tokenString})
	})

	// Protected routes with JWT middleware
	protected := e.Group("/secure")

	protected.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte("my-secret"),
	}))

	protected.GET("/profile", func(c echo.Context) error {
		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "invalid token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "invalid claims",
			})
		}

		username, ok := claims["username"].(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "username not found",
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"message": "Hello, " + username,
		})
	})

	e.Logger.Fatal(e.Start(":8080"))
}
