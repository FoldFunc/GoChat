package components

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/FoldFunc/GoChat/server/internal/app"
	"github.com/FoldFunc/GoChat/server/db"
)

func Login(w http.ResponseWriter, r *http.Request) {
	log.Println("/login called")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	var req app.LoginReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	userId, err := db.GetUserIdByNameDB(req.UserName)
	if err != nil {
		http.Error(w, "No such user", http.StatusBadRequest)
		return
	}
	password, err := db.GetUserPasswordByIdDB(userId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = db.SetUserLoggedInDB(userId, false)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if password != req.UserPassword {
		log.Printf("NO WAY:\npassword: %s;req.Password: %s", password, req.UserPassword)
		http.Error(w, "Invalid password", http.StatusForbidden)
		return
	}
	sessionId := app.GenerateId()
	app.SessionsMu.Lock()
	app.Sessions[strconv.Itoa(sessionId)] = userId
	app.SessionsMu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name: "session_id",
		Value: strconv.Itoa(sessionId),
		Path: "/",
		HttpOnly: true,
		Secure: false,
		SameSite: http.SameSiteLaxMode,
		MaxAge: 3600*24,
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "logged in",
	})
}
func Logout(w http.ResponseWriter, r *http.Request) {
	log.Println("/logout called")
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie("session_id")
	if err == nil {
		app.SessionsMu.Lock()
		delete(app.Sessions, cookie.Value)
		app.SessionsMu.Unlock()
		userId := r.Context().Value("userID").(int)
		err = db.SetUserLoggedInDB(userId, false)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "logged out",
	})
}

