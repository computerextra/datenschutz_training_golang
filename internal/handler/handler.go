package handler

import (
	"log/slog"
	"net/http"

	"computerextra/datenschutz_training_golang/db"

	"github.com/a-h/templ"
)

type Handler struct {
	logger   *slog.Logger
	database *db.PrismaClient
}

func New(logger *slog.Logger, database *db.PrismaClient) *Handler {
	return &Handler{
		logger:   logger,
		database: database,
	}
}

func Component(comp templ.Component) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		comp.Render(r.Context(), w)
	})
}
