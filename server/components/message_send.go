package components

import (
	"encoding/json"
	"net/http"
	"strconv"
	"log"

	"github.com/FoldFunc/GoChat/server/app"
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
	if !app.RoomExsists(req.RoomId) {
		http.Error(w, "No such room", http.StatusNotFound)
		return 
	}
	if !app.RoomPublic(req.RoomId) {
		http.Error(w, "Room not public", http.StatusForbidden)
		return
	}
	id := app.GenerateId()
	userId := r.Context().Value("userID").(int)
	message := app.Message{
		Id: id,
		UserId: userId,
		Body: req.Body,
	}
	for i := range app.R.Rooms {
		if app.R.Rooms[i].Id == req.RoomId {
			app.R.Rooms[i].Messages = append(app.R.Rooms[i].Messages, message)
		}
	}
	idInt := strconv.Itoa(id)
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
	if !app.RoomExsists(req.RoomId) {
		http.Error(w, "No such room", http.StatusNotFound)
		return 
	}
	if app.RoomPublic(req.RoomId) {
		http.Error(w, "Room public", http.StatusForbidden)
		return
	}
	if !app.UserInRoom(userId, req.RoomId) {
		http.Error(w, "You don't belong to this room", http.StatusForbidden)
		return
	}
	id := app.GenerateId()
	message := app.Message{
		Id: id,
		UserId: userId,
		Body: req.Body,
	}
	for i := range app.R.Rooms {
		if app.R.Rooms[i].Id == req.RoomId {
			app.R.Rooms[i].Messages = append(app.R.Rooms[i].Messages, message)
		}
	}
	idInt := strconv.Itoa(id)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "message created",
		"id": idInt,
	})
}
