package components

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/FoldFunc/GoChat/server/app"
	"github.com/FoldFunc/GoChat/server/db"
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
		Password: req.Password,
	}
	err = db.CreateUser(NewUser)
	if err != nil {
		http.Error(w, "Error while adding to the database", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"id": newId,
	})

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
	var roomType app.Type
	if req.RoomType {
		roomType = app.Type(app.TypePublic)
	} else {
		roomType = app.Type(app.TypePrivate)
	}
	newId := app.GenerateId()
	userId := r.Context().Value("userID").(int)
	currentUser, err := app.GetUserById(userId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	var admins []*app.User
	var users  []*app.User
	admins = append(admins, currentUser)
	users = append(users, currentUser)
	NewRoom := app.Room{
		Id: newId,
		UserId: userId,
		Name: req.RoomName,
		Type: roomType,
		Admins: admins,
		Users: users,
	}
	err = db.CreateRoom(NewRoom)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"id": newId,
	})
}
