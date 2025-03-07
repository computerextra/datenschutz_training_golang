package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func NeedAuth(store *sessions.CookieStore, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !hasAcceptedCookies(r) {
			http.Error(w, "Cookies müssen akzeptiert werden", http.StatusForbidden)
			return
		}

		session, _ := store.Get(r, "session-name")
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			// TODO: Show Login Page
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func hasAcceptedCookies(r *http.Request) bool {
	cookie, err := r.Cookie("cookies_accepted")
	if err != nil {
		return false
	}
	return cookie.Value == "true"
}
