package routes

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func Session(e *echo.Echo) {
	e.Use(
		session.Middleware(
			sessions.NewCookieStore([]byte("secret")),
		),
	)

	// Create Session
	e.GET("/create-session", func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil {
			return err
		}
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
		}
		sess.Values["foo"] = "bar"
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			return err
		}
		return c.NoContent(http.StatusOK)
	})

	// Read Session
	e.GET("/read-session", func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil {
			return err
		}
		return c.String(http.StatusOK, fmt.Sprintf("foo=%v\n", sess.Values["foo"]))
	})
}
