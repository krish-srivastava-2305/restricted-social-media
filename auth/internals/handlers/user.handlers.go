package handlers

import (
	"net/http"
	"time"

	"github.com/krish-srivastava-2305/internals/services"
	"github.com/labstack/echo/v4"
)

type RegisterRequest struct {
	Email            string    `json:"email"`
	Password         string    `json:"password"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	EmergencyContact string    `json:"emergency_contact"`
	DateOfBirth      time.Time `json:"date_of_birth"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register a new user
func RegisterHandler(c echo.Context) error {
	var req RegisterRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request format"})
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" || req.EmergencyContact == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "All fields are required"})
	}

	// Hash the password
	hashedPassword, err := services.HashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error hashing password"})
	}

	surveyDate := time.Now().Add(time.Hour * 24 * 7)

	// Register user
	user, err := services.RegisterUser(req.Email, hashedPassword, req.FirstName, req.LastName, req.EmergencyContact, req.DateOfBirth, surveyDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// Generate JWT token
	token, err := services.GenerateToken(user.Email, user.EmergencyContact, user.FirstName, user.SurveyDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error generating token"})
	}

	// Set cookie with the token
	cookie := http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(&cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "User registered successfully",
		"user":    user,
		"token":   token,
	})
}

// Login a user
func LoginHandler(c echo.Context) error {
	var req LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request format"})
	}

	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "All fields are required"})
	}

	// Authenticate user
	user, err := services.LoginUser(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	// Generate JWT token
	token, err := services.GenerateToken(user.Email, user.EmergencyContact, user.FirstName, user.SurveyDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error generating token"})
	}

	// Set cookie with the token
	cookie := http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(&cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Login successful",
		"user":    user,
		"token":   token,
	})
}

// Logout user
func LogoutHandler(c echo.Context) error {
	cookie := http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(&cookie)

	return c.JSON(http.StatusOK, echo.Map{"message": "Logged out successfully"})
}

func ProfileHandler(c echo.Context) error {
	email := c.Get("email").(string)

	user, err := services.GetUser(email)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "User found",
		"user":    user,
	})
}
