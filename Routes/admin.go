package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func AdminRoutes(e *echo.Echo) {
	g := e.Group("/admin")
	g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "joe" && password == "sescret" {
			return true, nil
		}
		return false, nil
	}))
	g.GET("/users", func(c echo.Context) error {
		return c.String(http.StatusOK, "Admin Users")
	})
}
