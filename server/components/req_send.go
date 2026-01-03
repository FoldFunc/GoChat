package components

import (
	"encoding/json"
	"net/http"

	"github.com/FoldFunc/GoChat/server/app"
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
	for i := range app.U.Users {
		if app.U.Users[i].Id == req.SendId {
			app.U.Users[i].ConnRequests = append(app.U.Users[i].ConnRequests, connRequest)
		}
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
	var requests []app.ConnReq
	for _, u := range app.U.Users {
		if u.Id == userId{
			requests = u.ConnRequests
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	err := json.NewEncoder(w).Encode(requests)
	if err != nil {
		http.Error(w, "Error while encoding json", http.StatusInternalServerError)
		return
	}
}
