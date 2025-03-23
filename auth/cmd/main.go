package main

import (
	"github.com/krish-srivastava-2305/config"
	"github.com/krish-srivastava-2305/internals/handlers"
	"github.com/krish-srivastava-2305/internals/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Connect to the database
	config.ConnectDB()
	defer config.CloseDB()

	e := echo.New()

	e.Use(middleware.CORS())

	e.POST("/api/auth/register", handlers.RegisterHandler)
	e.POST("/api/auth/login", handlers.LoginHandler)

	r := e.Group("/api/auth")

	r.Use(middlewares.AuthMiddleware)

	r.POST("/logout", handlers.LogoutHandler)
	r.GET("/profile", handlers.ProfileHandler)

	e.Logger.Fatal(e.Start(":8000"))
}
