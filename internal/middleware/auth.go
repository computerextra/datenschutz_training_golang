package middleware

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

func NeedAuth(store *sessions.CookieStore, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !hasAcceptedCookies(r) {
			// TODO: Show Error
			http.Error(w, "Cookies müssen akzeptiert werden", http.StatusForbidden)
			return
		}

		session, _ := store.Get(r, "session-name")
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			host := r.Host
			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}

			uri := fmt.Sprintf("%s://%s/signIn", scheme, host)
			http.Redirect(w, r, uri, http.StatusFound)
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

func IsAuthenticated(store *sessions.CookieStore, r *http.Request) bool {
	session, _ := store.Get(r, "session-name")
	auth, ok := session.Values["authenticated"].(bool)
	return ok && auth
}
