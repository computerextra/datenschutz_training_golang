package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session-name")
	// TODO: Check Credentials

	session.Values["authenticated"] = true
	session.Save(r, w)

	host := r.Host
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	uri := fmt.Sprintf("%s://%s", scheme, host)
	http.Redirect(w, r, uri, http.StatusFound)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session-name")
	session.Values["authenticated"] = false
	session.Save(r, w)
	host := r.Host
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	uri := fmt.Sprintf("%s://%s", scheme, host)
	http.Redirect(w, r, uri, http.StatusFound)
}
