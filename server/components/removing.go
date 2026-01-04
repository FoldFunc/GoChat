package components

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/FoldFunc/GoChat/server/app"
	"github.com/FoldFunc/GoChat/server/db"
)
func RemoveMessage(w http.ResponseWriter, r *http.Request) {
	log.Println("/removeMessage called")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	var req app.RemovemesReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("JSON ERROR: ", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	userId := r.Context().Value("userID").(int)
	if !app.UserExsists(userId) {
		http.Error(w, "User verification failed", http.StatusForbidden)
		return
	}
	if !app.RoomExsists(req.RoomId) {
		http.Error(w, "No such room", http.StatusBadRequest)
		return
	}
	if !app.MessageExsists(req.RoomId, req.MessId, userId) {
		http.Error(w, "No such message", http.StatusBadRequest)
		return
	}
	user, err := app.GetUserById(userId)
	if err != nil {
		http.Error(w, "No such user", http.StatusBadRequest)
		return
	}
	room, err := app.GetRoomById(req.RoomId)
	if err != nil {
		http.Error(w, "No such room", http.StatusBadRequest)
		return
	}
	err = db.RemoveMessage(*user, *room, req.MessId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "message deleted",
	}) 
}
func RemoveRoom(w http.ResponseWriter, r *http.Request) {
	log.Println("/removeRoom called")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	var req app.RemoveRoomReq 
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
	if !app.RoomExsistsToDelete(req.RoomId, userId) {
		http.Error(w, "You don't own the room", http.StatusForbidden)
		return
	}
	room, err := app.GetRoomById(req.RoomId)
	if err != nil {
		http.Error(w, "No such room", http.StatusBadRequest)
		return
	}
	err = db.RemoveRoom(*room)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "message deleted",
	}) 
}
