package app

import (
	"computerextra/datenschutz_training_golang/internal/component"
	"computerextra/datenschutz_training_golang/internal/handler"
	"computerextra/datenschutz_training_golang/internal/middleware"

	"fmt"
	"io/fs"
	"net/http"
	"os"
)

func (a *App) LoadPages(router *http.ServeMux) {
	h := handler.New(a.logger, a.database, a.ipresolver, a.store)

	router.Handle("GET /{$}", handler.ComponentWithContext(component.Index(), a.store))

	// Auth
	router.Handle("GET /signIn", handler.Component(component.SignIn()))
	router.HandleFunc("POST /signIn", h.Login)
	router.Handle("GET /signUp", handler.Component(component.SignUp()))
	router.HandleFunc("POST /signUp", h.Register)
	router.Handle("GET /signOut", a.auth(h.Logout))
	router.Handle("GET /loggedOut", handler.Component(component.LogOut()))
	router.HandleFunc("GET /verify", h.Verify)

	// Protected
	router.Handle("GET /prot", a.auth(h.Test))

	// Catch the Rest
	router.Handle("GET /", handler.Component(component.NotFound()))
}

func (a *App) auth(next func(http.ResponseWriter, *http.Request)) http.Handler {
	return middleware.NeedAuth(a.store, http.HandlerFunc(next))
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
