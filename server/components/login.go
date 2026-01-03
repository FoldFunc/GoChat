package components

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/FoldFunc/GoChat/server/app"
)

func Login(w http.ResponseWriter, r *http.Request) {
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
	user, err := app.GetUserByName(req.UserName)
	if err != nil {
		http.Error(w, "No such user", http.StatusBadRequest)
		return
	}
	if user.Password != req.UserPassword {
		http.Error(w, "Invalid password", http.StatusForbidden)
		return
	}
	userId := user.Id
	sessionId := app.GenerateId()
	app.Sessions[strconv.Itoa(sessionId)] = userId

	http.SetCookie(w, &http.Cookie{
		Name: "session_id",
		Value: strconv.Itoa(sessionId),
		Path: "/",
		HttpOnly: true,
		Secure: false,
		SameSite: http.SameSiteLaxMode,
		MaxAge: 3600*24,
	})
	json.NewEncoder(w).Encode(map[string]string{
		"message": "logged in",
	})
}
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("session_id")
    if err == nil {
        delete(app.Sessions, cookie.Value)
    }

    http.SetCookie(w, &http.Cookie{
        Name:   "session_id",
        Value:  "",
        Path:   "/",
        MaxAge: -1,
    })

    w.Write([]byte("logged out"))
}

