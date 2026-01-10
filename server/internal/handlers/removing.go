package components

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/FoldFunc/GoChat/server/internal/app"
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
	exsists, err := db.UserExists(userId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !exsists {
		http.Error(w, "User verification failed", http.StatusForbidden)
		return
	}
	roomExsists, err := db.RoomExistsDB(req.RoomId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !roomExsists {
		http.Error(w, "No such room", http.StatusBadRequest)
		return
	}
	messageExsists, err := db.MessageExistsDB(req.RoomId, req.MessId, userId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if ! messageExsists{
		http.Error(w, "No such message", http.StatusBadRequest)
		return
	}
	user, err := db.GetUserByIdDB(userId)
	if err != nil {
		http.Error(w, "No such user", http.StatusBadRequest)
		return
	}
	room, err := db.GetRoomByIDDB(req.RoomId)
	if err != nil {
		http.Error(w, "No such room", http.StatusBadRequest)
		return
	}
	err = db.RemoveMessage(user, room, req.MessId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
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
	exsists, err := db.UserExists(userId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !exsists {
		http.Error(w, "User verification failed", http.StatusForbidden)
		return
	}
	roomExsists, err := db.RoomExistsDB(req.RoomId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !roomExsists {
		http.Error(w, "No such room", http.StatusBadRequest)
		return
	}
	room, err := db.GetRoomByIDDB(req.RoomId)
	if err != nil {
		http.Error(w, "No such room", http.StatusBadRequest)
		return
	}
	err = db.RemoveRoom(room)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "room deleted",
	}) 
}
