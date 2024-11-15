package main

import (
	"net/http"

	"github.com/computerextra/datenschutz_training_golang/api"
	"github.com/computerextra/datenschutz_training_golang/db"
	"github.com/computerextra/datenschutz_training_golang/web"
	"github.com/computerextra/datenschutz_training_golang/web/template"
	"github.com/gin-gonic/gin"
)

func main() {
	err := db.Open()

	if err != nil {
		println("Error:", err)
	}

	// adding Admin, ignore error
	db.AddOrUpdateLogin(&db.Login{
		User:     "admin",
		Password: "password1",
		Permission: db.Permissions{
			AccessAdmin: db.Write,
			AccessBook:  db.Write,
		},
	})

	router := gin.Default()
	router.HTMLRender = &template.TemplRender{}

	// Serve Files
	router.Static("/static", "./static")
	router.StaticFile("/favicon.ico", "./static/favicon.ico")

	// web
	router.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "", template.Page()) })
	router.POST("/search", web.Search)
	router.GET("/book/add", web.ShowAddBook)
	router.GET("/book/:id", web.ShowEditBook)
	router.POST("/book", web.AddBook)
	router.PUT("/book/:id", web.EditBook)
	router.DELETE("/book/:id", web.DeleteBook)

	// api
	router.GET("/api/books", api.AuthMiddleware(api.BookReadOnly), api.GetBooks)
	router.GET("/api/book/:id", api.AuthMiddleware(api.BookReadOnly), api.GetBookById)
	router.POST("/api/book", api.AuthMiddleware(api.BookWrite), api.PostBook)
	router.DELETE("/api/book/:id", api.AuthMiddleware(api.BookWrite), api.DeleteBook)

	router.POST("/api/login", api.AuthMiddleware(api.LoginWrite), api.PostLogin)

	router.Run("localhost:8080")
}
