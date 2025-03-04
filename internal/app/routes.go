package app

import (
	"computerextra/datenschutz_training_golang/internal/component"
	"computerextra/datenschutz_training_golang/internal/handler"

	"fmt"
	"io/fs"
	"net/http"
	"os"
)

func (a *App) LoadPages(router *http.ServeMux) {
	// h := handler.New(a.logger, a.database)

	router.Handle("GET /{$}", handler.Component(component.Index()))

	// Catch the Rest
	router.Handle("GET /", handler.Component(component.NotFound()))
}

func (a *App) loadStaticFiles() (http.Handler, error) {
	if os.Getenv("BUILD_MODE") == "develop" {
		return http.FileServer(http.Dir("./static")), nil
	}

	static, err := fs.Sub(a.files, "static")
	if err != nil {
		return nil, fmt.Errorf("failed to subdir static: %w", err)
	}

	return http.FileServerFS(static), nil
}

func (a *App) loadRoutes() (http.Handler, error) {
	static, err := a.loadStaticFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to load static files: %w", err)
	}

	// Create new router
	router := http.NewServeMux()

	// this is the static file server
	router.Handle("GET /static/", http.StripPrefix("/static", static))

	a.LoadPages(router)

	return router, nil
}
