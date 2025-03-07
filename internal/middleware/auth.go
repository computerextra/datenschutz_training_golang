package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func NeedAuth(store *sessions.CookieStore, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			// TODO: Send Login Page
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// TODO: NYI

// func loginHandler(w http.ResponseWriter, r *http.Request) {
//     session, _ := store.Get(r, "session-name")
//     session.Values["authenticated"] = true
//     session.Save(r, w)
//     http.Redirect(w, r, "/", http.StatusFound)
// }

// func logoutHandler(w http.ResponseWriter, r *http.Request) {
//     session, _ := store.Get(r, "session-name")
//     session.Values["authenticated"] = false
//     session.Save(r, w)
//     http.Redirect(w, r, "/", http.StatusFound)
// }

// func hasAcceptedCookies(r *http.Request) bool {
//     cookie, err := r.Cookie("cookies_accepted")
//     if err != nil {
//         return false
//     }
//     return cookie.Value == "true"
// }
