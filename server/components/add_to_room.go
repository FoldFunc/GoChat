package components

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/FoldFunc/GoChat/server/app"
	"github.com/FoldFunc/GoChat/server/db"
)
func AddToCloseRoom(w http.ResponseWriter, r *http.Request) {
	log.Println("/addToCloseRoom called")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	var req app.AddToCloseRoomReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return 
	}
	adminID := r.Context().Value("userID").(int)
	exsists, err := db.UserExists(req.UserId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !exsists{
		http.Error(w, "User verification failed", http.StatusNotFound)
		return
	}
	isAdmin, err := db.IsUserAdminInRoomDB(adminID, req.RoomId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !isAdmin {
		http.Error(w, "AdminId is not an admin", http.StatusForbidden)
		return
	}
	currentUser, err := db.GetUserByIdDB(req.UserId)
	if err != nil {
		http.Error(w, "User not found", http.StatusForbidden)
		return
	}
	currentRoom, err := db.GetRoomByIDDB(req.RoomId)
	if err != nil {
		http.Error(w, "Room not found", http.StatusBadRequest)
		return
	}
	err = db.InsertUserCloseRoom(currentUser, currentRoom)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User added to a close room",
	})
}
func AddToOpenRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	var req app.AccesRoomReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
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
		http.Error(w, "User verification failed", http.StatusForbidden)
		return
	}
	roomExsists, err := db.RoomExistsDB(req.RoomId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !roomExsists{
		http.Error(w, "Room does not exsist", http.StatusBadRequest)
		return
	}
	isPublic, err := db.IsRoomPublicDB(req.RoomId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !isPublic {
		http.Error(w, "Room is not public", http.StatusBadRequest)
		return
	}
	currentUser, err := db.GetUserByIdDB(userId) 
	if err != nil {
		http.Error(w, "No such user", http.StatusForbidden)
		return
	}
	currentRoom, err := db.GetRoomByIDDB(req.RoomId) 
	if err != nil {
		http.Error(w, "No such room", http.StatusForbidden)
		return
	}
	err = db.InsertUserCloseRoom(currentUser, currentRoom)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User added to a close room",
	})
}
