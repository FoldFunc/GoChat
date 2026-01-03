package components

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/FoldFunc/GoChat/server/app"
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
	if !app.UserExsists(req.UserId) {
		http.Error(w, "User verification failed", http.StatusNotFound)
		return
	}
	if !app.UserExsists(req.AdminId) {
		http.Error(w, "Admin verification failed", http.StatusNotFound)
		return
	}
	if !app.IsAdmin(req.AdminId, req.RoomId) {
		http.Error(w, "AdminId is not an admin", http.StatusForbidden)
		return
	}
	currentUser, err := app.GetUserById(req.UserId)
	if err != nil {
		http.Error(w, "User not found", http.StatusForbidden)
		return
	}
	for i := range app.R.Rooms {
		if app.R.Rooms[i].Id == req.RoomId {
			app.R.Rooms[i].Users = append(app.R.Rooms[i].Users, currentUser)
		}
	}
	currentRoom, err := app.GetRoomById(req.RoomId)
	if err != nil {
		http.Error(w, "Room not found", http.StatusBadRequest)
		return
	}
	for i := range app.U.Users {
		if app.U.Users[i].Id == req.UserId {
			app.U.Users[i].Rooms = append(app.U.Users[i].Rooms, currentRoom)
		}
	}
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
	if !app.UserExsists(req.UserId) {
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
	currentUser, err := app.GetUserById(req.UserId) 
	if err != nil {
		http.Error(w, "No such user", http.StatusForbidden)
		return
	}
	currentRoom, err := app.GetRoomById(req.RoomId) 
	if err != nil {
		http.Error(w, "No such room", http.StatusForbidden)
		return
	}
	for i := range app.R.Rooms {
		if app.R.Rooms[i].Id == req.RoomId {
			app.R.Rooms[i].Users = append(app.R.Rooms[i].Users, currentUser)
		}
	}
	for i := range app.U.Users {
		if app.U.Users[i].Id == req.UserId {
			app.U.Users[i].Rooms = append(app.U.Users[i].Rooms, currentRoom)
		}
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User added to the room",
	})
}
