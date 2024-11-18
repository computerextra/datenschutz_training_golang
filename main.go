package main

import (
	"io"
	"net/http"
	"os"

	routes "github.com/computerextra/datenschutz_training_golang/Routes"
	templates "github.com/computerextra/datenschutz_training_golang/Templates"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Link struct {
	Name string
	Url  string
}

// Handling Request
type User struct {
	Name  string `json:"name" xml:"name" form:"name" query:"name"`
	Email string `json:"email" xml:"email" form:"email" query:"email"`
}

// Cookies
// https://echo.labstack.com/docs/cookies

func main() {
	e := echo.New()

	// Auto TLS
	// e.AutoTLSManager.HostPolicy = autocert.HostWhitelist("johanneskirchner.net")
	// Cache Certificate to avaid issues with rate Limits (https://letsencrypt.org/docs/rate-limits)
	// e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")

	// Remove trailing slash
	e.Pre(middleware.RemoveTrailingSlash())

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Sessions
	routes.Session(e)

	routes.AdminRoutes(e)

	// Route Level Middleware
	track := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			println("request to /users")
			return next(c)
		}
	}
	e.GET("/users", func(c echo.Context) error {
		return c.String(http.StatusOK, "/users")
	}, track)

	// Serve Static Files
	// Serve any file from static directory for path /static/*.
	e.Static("/static", "static")

	e.GET("/", func(c echo.Context) error {

		return templates.Render(c, http.StatusOK, templates.Home())
	})
	// Query Parameters
	e.GET("/show", show)

	// Form application/x-www-form-urlencoded
	e.POST("/save", save)

	// Routing
	// Handling Request
	e.POST("/users", func(c echo.Context) error {
		u := new(User)
		if err := c.Bind(u); err != nil {
			return err
		}
		return c.JSON(http.StatusCreated, u)
		// or
		// return c.XML(http.StatusCreated, u)
	})
	e.GET("/users/:id", getUser)
	// e.PUT("/users/:id", updateUser)
	// e.DELETE("/users/:id", deleteUser)

	e.Logger.Fatal(e.Start(":8080"))
	// Auto TLS
	// e.Logger.Fatal(e.StartAutoTLS(":443"))
}

// Path Parameters
// e.GET("/users/:id", getUser)
func getUser(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}

// Query Parameters
func show(c echo.Context) error {
	// Get team and Member from the querystring
	team := c.QueryParam("team")
	member := c.QueryParam("member")
	return c.String(http.StatusOK, "team: "+team+", member:"+member)
}

// Form application/x-www-form-urlencoded
func save(c echo.Context) error {
	// Get name and email
	name := c.FormValue("name")
	email := c.FormValue("email")
	return c.String(http.StatusOK, "name: "+name+", email: "+email)
}

// Form multipart/form-data
func saveForm(c echo.Context) error {
	// Get name
	name := c.FormValue("name")
	// Get Avatar
	avatar, err := c.FormFile("avatar")
	if err != nil {
		return err
	}

	// Source
	src, err := avatar.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(avatar.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, "<b>Thank You! "+name+"</b>")

}
