package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"computerextra/datenschutz_training_golang/db"
	"computerextra/datenschutz_training_golang/internal/service/realip"

	"github.com/a-h/templ"
)

type Handler struct {
	logger     *slog.Logger
	database   *db.PrismaClient
	ipResolver *realip.Service
}

func New(logger *slog.Logger, database *db.PrismaClient, ipService *realip.Service) *Handler {
	return &Handler{
		logger:     logger,
		database:   database,
		ipResolver: ipService,
	}
}

func Component(comp templ.Component) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		comp.Render(r.Context(), w)
	})
}

func (h *Handler) Test(w http.ResponseWriter, r *http.Request) {

	host := r.Host
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	uri := fmt.Sprintf("%s://%s", scheme, host)

	http.Redirect(w, r, uri, http.StatusFound)
}
