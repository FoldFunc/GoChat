package components

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/FoldFunc/GoChat/server/app"
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
	if !app.UserExsists(req.UserId) {
		http.Error(w, "User verification failed", http.StatusForbidden)
		return
	}
	if !app.RoomExsists(req.RoomId) {
		http.Error(w, "No such room", http.StatusBadRequest)
		return
	}
	if !app.MessageExsists(req.RoomId, req.MessId, req.UserId) {
		http.Error(w, "No such message", http.StatusBadRequest)
		return
	}
	for i := range app.R.Rooms {
		if app.R.Rooms[i].Id == req.RoomId {
			for j := range app.R.Rooms[i].Messages {
				if app.R.Rooms[i].Messages[j].Id == req.MessId {
					app.R.Rooms[i].Messages = append(app.R.Rooms[i].Messages[:j], app.R.Rooms[i].Messages[j+1:]...)
					json.NewEncoder(w).Encode(map[string]string{
						"message": "message removed",
					})
					return
				}
			}
		}
	}
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
	if !app.UserExsists(req.UserId) {
		http.Error(w, "User verification failed", http.StatusForbidden)
		return
	}
	if !app.RoomExsistsToDelete(req.RoomId, req.UserId) {
		http.Error(w, "You don't own the room", http.StatusForbidden)
		return
	}
	for i := range app.R.Rooms {
		if app.R.Rooms[i].Id == req.RoomId {
			app.R.Rooms = append(app.R.Rooms[:i], app.R.Rooms[i+1:]...)

			json.NewEncoder(w).Encode(map[string]string{
				"message": "Room deleted",
			})
			return
		}
	}
}
