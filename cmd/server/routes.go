package main

import (
	"net/http"

	"github.com/ktnyt/GoWebStack/cmpt"
	"github.com/labstack/echo/v4"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(e *echo.Echo) {
	e.Static("/assets", "assets")

	e.GET("/", func(c echo.Context) error {
		component := cmpt.HTML(cmpt.Hello("nano"))
		return component.Render(c.Request().Context(), c.Response().Writer)
	})

	e.GET("/status", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	go http.Get("http://localhost:6641/reload")
}
