package components

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/FoldFunc/GoChat/server/internal/app"
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
	NewUser := app.UserData{
		Id: newId,
		Name: req.UserName,
		ConnType: req.ConnType,
	}
	err = db.CreateUser(NewUser, req.Password)
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
	currentUser, err := db.GetUserByIdDB(userId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	var admins []app.UserData
	var users  []app.UserData
	admins = append(admins, currentUser)
	users = append(users, currentUser)
	NewRoom := app.RoomData{
		Id: newId,
		UserId: userId,
		Name: req.RoomName,
		Type: roomType,
	}
	err = db.CreateRoom(NewRoom)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = db.AddToRoomUserDB(NewRoom, userId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = db.AddUserAsAdminDB(NewRoom, userId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"id": newId,
	})
}
