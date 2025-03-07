package handler

import (
	"computerextra/datenschutz_training_golang/db"
	"computerextra/datenschutz_training_golang/internal/component"
	"computerextra/datenschutz_training_golang/internal/utils"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	mail := r.FormValue("mail")
	pass := r.FormValue("password")

	user, err := h.database.User.FindUnique(
		db.User.Email.Equals(mail),
	).Exec(r.Context())
	if err != nil {
		h.logger.Error("failed to find user", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !utils.CheckPasswordHash(pass, user.PasswordHash) {
		host := r.Host
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}

		uri := fmt.Sprintf("%s://%s/signIn", scheme, host)
		http.Redirect(w, r, uri, http.StatusFound)
		return
	}

	_, ok := user.EmailVerified()
	if !ok {
		host := r.Host
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}

		uri := fmt.Sprintf("%s://%s/verify", scheme, host)
		http.Redirect(w, r, uri, http.StatusFound)
		return
	}

	session, _ := h.store.Get(r, "session-name")
	session.Values["authenticated"] = true
	session.Values["user"] = user.Name
	session.Values["mail"] = user.Email
	session.Values["admin"] = user.Admin
	session.Values["chef"] = user.Chef

	session.Save(r, w)
	host := r.Host
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	uri := fmt.Sprintf("%s://%s", scheme, host)
	http.Redirect(w, r, uri, http.StatusFound)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	mail := r.FormValue("mail")
	pass := r.FormValue("password")

	ctx := r.Context()

	if len(name) < 1 {
		h.logger.Error("name buggy", slog.Any("error", name))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(mail) < 1 {
		h.logger.Error("mail buggy", slog.Any("error", mail))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(pass) < 1 {
		h.logger.Error("pass buggy", slog.Any("error", pass))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := generateToken()
	if err != nil {
		h.logger.Error("failed to generate token", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	passwordHash, err := utils.HashPassword(pass)
	if err != nil {
		h.logger.Error("failed to hash password", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = h.database.User.CreateOne(
		db.User.Email.Set(mail),
		db.User.PasswordHash.Set(passwordHash),
		db.User.VerificationToken.Set(token),
		db.User.Name.Set(name),
	).Exec(ctx)
	if err != nil {
		h.logger.Error("failed to generate user", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	host := r.Host
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	uri := fmt.Sprintf("%s://%s", scheme, host)
	err = utils.SendVerificationMail(mail, token, uri)
	if err != nil {
		h.logger.Error("failed to send mail", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	uri = fmt.Sprintf("%s://%s/verify", scheme, host)
	http.Redirect(w, r, uri, http.StatusFound)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session-name")
	session.Values["authenticated"] = false
	session.Values["user"] = nil
	session.Values["mail"] = nil
	session.Values["admin"] = nil
	session.Values["chef"] = nil
	session.Save(r, w)
	host := r.Host
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	uri := fmt.Sprintf("%s://%s/loggedOut", scheme, host)
	http.Redirect(w, r, uri, http.StatusFound)
}

func (h *Handler) Verify(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	ctx := r.Context()

	if token == "" {
		component.Verify(false).Render(ctx, w)
		return
	}

	user, err := h.database.User.FindFirst(
		db.User.VerificationToken.Equals(token),
	).Exec(ctx)

	if err != nil {
		h.logger.Error("failed to get user", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user == nil {
		component.Verify(false).Render(ctx, w)
		return
	}

	if token == user.VerificationToken {
		_, err := h.database.User.FindUnique(
			db.User.ID.Equals(user.ID),
		).Update(
			db.User.EmailVerified.Set(time.Now()),
		).Exec(ctx)
		if err != nil {
			h.logger.Error("failed to set verification", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		component.Verify(true).Render(ctx, w)
		return
	} else {
		component.Verify(false).Render(ctx, w)
		return
	}
}
