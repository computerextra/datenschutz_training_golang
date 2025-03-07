package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"computerextra/datenschutz_training_golang/db"
	"computerextra/datenschutz_training_golang/internal/service/realip"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
)

type Handler struct {
	logger     *slog.Logger
	database   *db.PrismaClient
	ipResolver *realip.Service
	store      *sessions.CookieStore
}

func New(logger *slog.Logger, database *db.PrismaClient, ipService *realip.Service, store *sessions.CookieStore) *Handler {
	return &Handler{
		logger:     logger,
		database:   database,
		ipResolver: ipService,
		store:      store,
	}
}

func Component(comp templ.Component) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		comp.Render(r.Context(), w)
	})
}

type User struct {
	Name   string
	Mail   string
	Admin  bool
	Chef   bool
	Authed bool
}

type contextKey string

func setUser(ctx context.Context, u *User) context.Context {
	var Name contextKey = "name"
	var Mail contextKey = "mail"
	var Admin contextKey = "admin"
	var Chef contextKey = "chef"
	var Authed contextKey = "auth"

	ctx = context.WithValue(ctx, Name, u.Name)
	ctx = context.WithValue(ctx, Mail, u.Mail)
	ctx = context.WithValue(ctx, Admin, u.Admin)
	ctx = context.WithValue(ctx, Chef, u.Chef)
	ctx = context.WithValue(ctx, Authed, u.Authed)

	return ctx
}

func ComponentWithContext(comp templ.Component, store *sessions.CookieStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		nameInterface := session.Values["user"]
		authenticatedInterface := session.Values["authenticated"]
		mailInterface := session.Values["mail"]
		adminInterface := session.Values["admin"]
		chefInterface := session.Values["chef"]

		var name string = ""
		var mail string = ""
		var admin bool = false
		var chef bool = false
		var authenticated bool = false

		if nameInterface != nil {
			name, _ = nameInterface.(string)
		}
		if mailInterface != nil {
			mail, _ = mailInterface.(string)
		}
		if adminInterface != nil {
			admin, _ = adminInterface.(bool)
		}
		if chefInterface != nil {
			chef, _ = chefInterface.(bool)
		}
		if authenticatedInterface != nil {
			authenticated, _ = authenticatedInterface.(bool)
		}

		ctx := setUser(r.Context(), &User{
			Name:   name,
			Mail:   mail,
			Admin:  admin,
			Chef:   chef,
			Authed: authenticated,
		})

		w.Header().Add("Content-Type", "text/html")
		comp.Render(ctx, w)
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
