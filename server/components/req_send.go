package components

import (
	"encoding/json"
	"net/http"

	"github.com/FoldFunc/GoChat/server/app"
	"github.com/FoldFunc/GoChat/server/db"
)
func SendUserRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var req app.SendUserReq
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	userId := r.Context().Value("userID").(int)
	if !app.UserExsists(req.SendId) {
		http.Error(w, "No such user", http.StatusNotFound)
		return
	}
	if !app.UserPrivate(req.SendId) {
		http.Error(w, "No need to send the reques, user public", http.StatusNotAcceptable)
		return
	}
	connRequest := app.ConnReq{
		FromReqId: userId,
		Message: req.Message,
	}
	err = db.AddUserReq(connRequest, req.SendId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User request sent",
	})
}
func ViewUserRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request", http.StatusMethodNotAllowed)
		return
	}
	userId := r.Context().Value("userID").(int)
	if !app.UserPrivate(userId) {
		http.Error(w, "No need for this method, user public", http.StatusNotAcceptable)
		return
	}
	user, err := app.GetUserById(userId)
	if err != nil {
		http.Error(w, "No such user", http.StatusBadRequest)
		return 
	}
	requests, err := db.GetConnReq(*user)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return 
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	err = json.NewEncoder(w).Encode(requests)
	if err != nil {
		http.Error(w, "Error while encoding json", http.StatusInternalServerError)
		return
	}
}
