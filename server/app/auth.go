package app

import (
	"context"
	"log"
	"net/http"
)

// When security: Cookie doesn't expire.
func AuthCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Called auth: ", r.URL.Path)
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Invalid cookie", http.StatusUnauthorized)
			return
		}
		SessionsMu.RLock()
		userID, ok := Sessions[cookie.Value]
		SessionsMu.RUnlock()
		if !ok {
			http.Error(w, "unathorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
