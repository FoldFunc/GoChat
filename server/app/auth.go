package app

import (
	"context"
	"net/http"
)

func AuthCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Invalid cookie", http.StatusUnauthorized)
			return
		}
		userID, ok := Sessions[cookie.Value]
		if !ok {
			http.Error(w, "unathorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
