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
	if !app.UserExsists(req.UserId) {
		http.Error(w, "User verification failed", http.StatusNotFound)
		return
	}
	if !app.IsAdmin(adminID, req.RoomId) {
		http.Error(w, "AdminId is not an admin", http.StatusForbidden)
		return
	}
	currentUser, err := app.GetUserById(req.UserId)
	if err != nil {
		http.Error(w, "User not found", http.StatusForbidden)
		return
	}
	currentRoom, err := app.GetRoomById(req.RoomId)
	if err != nil {
		http.Error(w, "Room not found", http.StatusBadRequest)
		return
	}
	err = db.InsertUserCloseRoom(*currentUser, *currentRoom)
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
	if !app.UserExsists(userId) {
		http.Error(w, "User verification failed", http.StatusForbidden)
		return
	}
	if !app.RoomExsists(req.RoomId) {
		http.Error(w, "Room does not exsist", http.StatusBadRequest)
		return
	}
	if !app.RoomPublic(req.RoomId) {
		http.Error(w, "Room is not public", http.StatusBadRequest)
		return
	}
	currentUser, err := app.GetUserById(userId) 
	if err != nil {
		http.Error(w, "No such user", http.StatusForbidden)
		return
	}
	currentRoom, err := app.GetRoomById(req.RoomId) 
	if err != nil {
		http.Error(w, "No such room", http.StatusForbidden)
		return
	}
	err = db.InsertUserCloseRoom(*currentUser, *currentRoom)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User added to a close room",
	})
}
