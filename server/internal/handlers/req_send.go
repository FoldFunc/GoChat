package components

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/FoldFunc/GoChat/server/internal/app"
	"github.com/FoldFunc/GoChat/server/db"
)
func SendUserRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("/sendUserRequest called")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var req app.SendUserReq
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		log.Println("JSON ERROR: ", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	userId := r.Context().Value("userID").(int)
	exsists, err := db.UserExists(userId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !exsists{
		http.Error(w, "No such user", http.StatusNotFound)
		return
	}
	isPrivate, err := db.IsUserPrivateDB(req.SendId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !isPrivate {
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
	w.Header().Set("Content-Type", "application/json")
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
	isPrivate, err := db.IsUserPrivateDB(userId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !isPrivate {
		http.Error(w, "No need to send the reques, user public", http.StatusNotAcceptable)
		return
	}
	user, err := db.GetUserByIdDB(userId)
	if err != nil {
		http.Error(w, "No such user", http.StatusBadRequest)
		return 
	}
	requests, err := db.GetConnReq(user)
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
func ViewUserRequestsFromUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request", http.StatusMethodNotAllowed)
		return
	}
	var req app.ViewUserRequestsFromUserReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	userId := r.Context().Value("userID").(int)
	isPrivate, err := db.IsUserPrivateDB(userId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !isPrivate {
		http.Error(w, "No need to send the reques, user public", http.StatusNotAcceptable)
		return
	}
	user, err := db.GetUserByIdDB(userId)
	if err != nil {
		http.Error(w, "No such user", http.StatusBadRequest)
		return 
	}
	userFrom, err := db.GetUserByIdDB(req.UserId)
	if err != nil {
		http.Error(w, "No such user", http.StatusBadRequest)
		return 
	}
	requests, err := db.GetConnReqFrom(user, userFrom)
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
func AcceptRequestFromUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	var req app.AcceptRequestFromUserReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invlaid request", http.StatusBadRequest)
		return
	}
	userId := r.Context().Value("userID").(int)
	exsists, err := db.RequestExists(req.RequestId, userId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !exsists {
		http.Error(w, "Invalid request id", http.StatusBadRequest)
		return
	}
	err = db.RequestAccept(userId, req.RequestId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "request accepted",
	})
}
