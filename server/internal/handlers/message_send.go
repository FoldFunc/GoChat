package components

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/FoldFunc/GoChat/server/internal/app"
	"github.com/FoldFunc/GoChat/server/db"
)
func SendMessageOpenRoom(w http.ResponseWriter, r *http.Request) {
	log.Println("/sendMessageOpenRoom called")
	if r.Method != http.MethodPost {
		http.Error(w, "Only get requests allowed", http.StatusMethodNotAllowed)
		return
	}
	var req app.SendMessreq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	log.Printf("room_id: %d;message: %s", req.RoomId, req.Body)
	exsists, err := db.RoomExistsDB(req.RoomId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !exsists {
		http.Error(w, "No such room", http.StatusBadRequest)
		return
	}
	isPublic, err := db.IsRoomPublicDB(req.RoomId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !isPublic{
		http.Error(w, "Room is private", http.StatusBadRequest)
		return
	}
	id := app.GenerateId()
	userId := r.Context().Value("userID").(int)
	message := app.Message{
		Id: id,
		UserId: userId,
		Body: req.Body,
	}
	room, err := db.GetRoomByIDDB(req.RoomId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = db.InsertMessageRoom(message, room)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	idInt:= strconv.Itoa(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "message created",
		"id": idInt,
	})
}
func SendMessageCloseRoom(w http.ResponseWriter, r *http.Request) {
	log.Println("/sendMessageCloseRoom called")
	if r.Method != http.MethodPost {
		http.Error(w, "Only get requests allowed", http.StatusMethodNotAllowed)
		return
	}
	var req app.SendMessreq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	userId := r.Context().Value("userID").(int)
	exsists, err := db.RoomExistsDB(req.RoomId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !exsists {
		http.Error(w, "No such room", http.StatusBadRequest)
		return
	}
	isPublic, err := db.IsRoomPublicDB(req.RoomId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if isPublic {
		http.Error(w, "Room public", http.StatusBadRequest)
		return
	}
	inRoom, err := db.IsUserInRoomDB(userId, req.RoomId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !inRoom {
		http.Error(w, "You don't belong to this room", http.StatusForbidden)
		return
	}
	id := app.GenerateId()
	message := app.Message{
		Id: id,
		UserId: userId,
		Body: req.Body,
	}
	room, err  := db.GetRoomByIDDB(req.RoomId)
	if err != nil {
		http.Error(w, "No such room", http.StatusBadRequest)
		return
	}
	err = db.InsertMessageRoom(message, room)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	idInt := strconv.Itoa(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "message created",
		"id": idInt,
	})
}
