package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

var db *sql.DB

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {

	// Connect to MySQL
	var err error
	dsn := "root:" + "Carkey2003@" + "@tcp(localhost:3306)/simple_auth"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	error := db.Ping()
	if error != nil {
		log.Fatal("Failed to ping database:", error)
	}
	// Create Echo
	e := echo.New()

	// Routes
	e.POST("/register", register)
	e.POST("/login", login)

	e.Logger.Fatal(e.Start(":8080"))
}

// Register
func register(c echo.Context) error {
	user := new(User)

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid data")
	}

	_, err := db.Exec("INSERT INTO users(email, password) VALUES(?, ?)", user.Email, user.Password)
	if err != nil {
		log.Println("DB Insert Error:", err) // Log actual DB error here
		return c.JSON(http.StatusInternalServerError, "Error saving user")
	}

	return c.JSON(http.StatusOK, "User registered")
}

// Login
func login(c echo.Context) error {
	user := new(User)

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid data")
	}

	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE email = ?", user.Email).
		Scan(&storedPassword)

	if err != nil {
		log.Println("DB Query Error:", err) // Log actual DB error here
		return c.JSON(http.StatusUnauthorized, "User not found")
	}

	if storedPassword != user.Password {
		return c.JSON(http.StatusUnauthorized, "Wrong password")
	}

	return c.JSON(http.StatusOK, "Login successful")
}
