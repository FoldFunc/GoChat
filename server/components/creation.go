package components

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/FoldFunc/GoChat/server/app"
)
func NewUser(w http.ResponseWriter, r *http.Request) {
	log.Println("/newUser called")
	if r.Method != http.MethodPost {
		http.Error(w, "Only post requests allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var req app.NewUserReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	newId := app.GenerateId()
	var connType app.Type
	if req.ConnType {
		connType = app.TypePublic
	} else {
		connType = app.TypePrivate
	}
	NewUser := app.User{
		Id: newId,
		Name: req.UserName,
		ConnType: connType,
	}
	app.U.Users = append(app.U.Users, NewUser)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"id": newId,
	})
	for _, u := range app.U.Users {
		fmt.Printf("User id: %d; UserName: %s\n", u.Id, u.Name)
	}

}
func NewRoom(w http.ResponseWriter, r *http.Request) {
	log.Println("/newRoom called")
	if r.Method != http.MethodPost {
		http.Error(w, "Only post requests allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var req app.NewRoomReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if !app.UserExsists(req.UserId) {
		http.Error(w, "User validation failed", http.StatusForbidden)
		return
	}
	var roomType app.Type
	if req.RoomType {
		roomType = app.Type(app.TypePublic)
	} else {
		roomType = app.Type(app.TypePrivate)
	}
	newId := app.GenerateId()
	currentUser, err := app.GetUserById(req.UserId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	var admins []app.User
	var users  []app.User
	admins = append(admins, currentUser)
	users = append(users, currentUser)
	NewRoom := app.Room{
		Id: newId,
		UserId: req.UserId,
		Name: req.RoomName,
		Type: roomType,
		Admins: admins,
		Users: users,
	}
	app.R.Rooms = append(app.R.Rooms, NewRoom)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"id": newId,
	})
}
