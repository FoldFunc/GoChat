package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

)
func hello(w http.ResponseWriter, r *http.Request) {
	log.Println("/ handler called")
}
func newUser(w http.ResponseWriter, r *http.Request) {
	log.Println("/newUser called")
	if r.Method != http.MethodPost {
		http.Error(w, "Only post requests allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var req NewUserReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	newId := generateId()
	NewUser := User{
		Id: newId,
		Name: req.UserName,
	}
	U.Users = append(U.Users, NewUser)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"id": newId,
	})
	for _, u := range U.Users {
		fmt.Printf("User id: %d; UserName: %s\n", u.Id, u.Name)
	}

}
func newRoom(w http.ResponseWriter, r *http.Request) {
	log.Println("/newRoom called")
	if r.Method != http.MethodPost {
		http.Error(w, "Only post requests allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var req NewRoomReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if !userExsists(req.UserId) {
		http.Error(w, "User validation failed", http.StatusForbidden)
		return
	}
	var roomType Type
	if req.RoomType {
		roomType = Type(TypePublic)
	} else {
		roomType = Type(TypePrivate)
	}
	newId := generateId()
	currentUser, err := getUserById(req.UserId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	var admins []User
	var users  []User
	admins = append(admins, currentUser)
	users = append(users, currentUser)
	NewRoom := Room{
		Id: newId,
		UserId: req.UserId,
		Name: req.RoomName,
		Type: roomType,
		Admins: admins,
		Users: users,
	}
	R.Rooms = append(R.Rooms, NewRoom)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"id": newId,
	})
}
func sendMessageCloseRoom(w http.ResponseWriter, r *http.Request) {
	log.Println("/sendMessageCloseRoom called")
	if r.Method != http.MethodPost {
		http.Error(w, "Only get requests allowed", http.StatusMethodNotAllowed)
		return
	}
	var req SendMessreq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if !userExsists(req.UserId) {
		http.Error(w, "User verification failed", http.StatusForbidden)
		return
	}
	if !roomExsists(req.RoomId) {
		http.Error(w, "No such room", http.StatusNotFound)
		return 
	}
	if roomPublic(req.RoomId) {
		http.Error(w, "Room public", http.StatusForbidden)
		return
	}
	if !UserInRoom(req.UserId, req.RoomId) {
		http.Error(w, "You don't belong to this room", http.StatusForbidden)
		return
	}
	id := generateId()
	message := Message{
		Id: id,
		UserId: req.UserId,
		Body: req.Body,
	}
	for i := range R.Rooms {
		if R.Rooms[i].Id == req.RoomId {
			R.Rooms[i].Messages = append(R.Rooms[i].Messages, message)
		}
	}
	idInt := strconv.Itoa(id)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "message created",
		"id": idInt,
	})
}
func sendMessageOpenRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only get requests allowed", http.StatusMethodNotAllowed)
		return
	}
	var req SendMessreq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if !userExsists(req.UserId) {
		http.Error(w, "User verification failed", http.StatusForbidden)
		return
	}
	if !roomExsists(req.RoomId) {
		http.Error(w, "No such room", http.StatusNotFound)
		return 
	}
	if !roomPublic(req.RoomId) {
		http.Error(w, "Room not public", http.StatusForbidden)
		return
	}
	id := generateId()
	message := Message{
		Id: id,
		UserId: req.UserId,
		Body: req.Body,
	}
	for i := range R.Rooms {
		if R.Rooms[i].Id == req.RoomId {
			R.Rooms[i].Messages = append(R.Rooms[i].Messages, message)
		}
	}
	idInt := strconv.Itoa(id)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "message created",
		"id": idInt,
	})
}
func removeMessage(w http.ResponseWriter, r *http.Request) {
	log.Println("/removeMessage called")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	var req RemovemesReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("JSON ERROR: ", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if !userExsists(req.UserId) {
		http.Error(w, "User verification failed", http.StatusForbidden)
		return
	}
	if !roomExsists(req.RoomId) {
		http.Error(w, "No such room", http.StatusBadRequest)
		return
	}
	if !messageExsists(req.RoomId, req.MessId, req.UserId) {
		http.Error(w, "No such message", http.StatusBadRequest)
		return
	}
	for i := range R.Rooms {
		if R.Rooms[i].Id == req.RoomId {
			for j := range R.Rooms[i].Messages {
				if R.Rooms[i].Messages[j].Id == req.MessId {
					R.Rooms[i].Messages = append(R.Rooms[i].Messages[:j], R.Rooms[i].Messages[j+1:]...)
					json.NewEncoder(w).Encode(map[string]string{
						"message": "message removed",
					})
					return
				}
			}
		}
	}
}
func removeRoom(w http.ResponseWriter, r *http.Request) {
	log.Println("/removeRoom called")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	var req RemoveRoomReq 
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if !userExsists(req.UserId) {
		http.Error(w, "User verification failed", http.StatusForbidden)
		return
	}
	if !roomExsistsToDelete(req.RoomId, req.UserId) {
		http.Error(w, "You don't own the room", http.StatusForbidden)
		return
	}
	for i := range R.Rooms {
		if R.Rooms[i].Id == req.RoomId {
			R.Rooms = append(R.Rooms[:i], R.Rooms[i+1:]...)

			json.NewEncoder(w).Encode(map[string]string{
				"message": "Room deleted",
			})
			return
		}
	}
}
func addToCloseRoom(w http.ResponseWriter, r *http.Request) {
	log.Println("/addToCloseRoom called")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	var req AddToCloseRoomReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return 
	}
	if !userExsists(req.UserId) {
		http.Error(w, "User verification failed", http.StatusNotFound)
		return
	}
	if !userExsists(req.AdminId) {
		http.Error(w, "Admin verification failed", http.StatusNotFound)
		return
	}
	if !IsAdmin(req.AdminId, req.RoomId) {
		http.Error(w, "AdminId is not an admin", http.StatusForbidden)
		return
	}
	currentUser, err := getUserById(req.UserId)
	if err != nil {
		http.Error(w, "User not found", http.StatusForbidden)
		return
	}
	for i := range R.Rooms {
		if R.Rooms[i].Id == req.RoomId {
			R.Rooms[i].Users = append(R.Rooms[i].Users, currentUser)
		}
	}
	currentRoom, err := getRoomById(req.RoomId)
	if err != nil {
		http.Error(w, "Room not found", http.StatusBadRequest)
		return
	}
	for i := range U.Users {
		if U.Users[i].Id == req.UserId {
			U.Users[i].Rooms = append(U.Users[i].Rooms, currentRoom)
		}
	}
}
func addToOpenRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	var req AccesRoomReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if !userExsists(req.UserId) {
		http.Error(w, "User verification failed", http.StatusForbidden)
		return
	}
	if !roomExsists(req.RoomId) {
		http.Error(w, "Room does not exsist", http.StatusBadRequest)
		return
	}
	if !roomPublic(req.RoomId) {
		http.Error(w, "Room is not public", http.StatusBadRequest)
		return
	}
	currentUser, err := getUserById(req.UserId) 
	if err != nil {
		http.Error(w, "No such user", http.StatusForbidden)
		return
	}
	currentRoom, err := getRoomById(req.RoomId) 
	if err != nil {
		http.Error(w, "No such room", http.StatusForbidden)
		return
	}
	for i := range R.Rooms {
		if R.Rooms[i].Id == req.RoomId {
			R.Rooms[i].Users = append(R.Rooms[i].Users, currentUser)
		}
	}
	for i := range U.Users {
		if U.Users[i].Id == req.UserId {
			U.Users[i].Rooms = append(U.Users[i].Rooms, currentRoom)
		}
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User added to the room",
	})
}
var R GlobalRooms
var U GlobalUsers
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	mux.HandleFunc("/newUser", newUser)
	mux.HandleFunc("/newRoom", newRoom)
	mux.HandleFunc("/addToOpenRoom",addToOpenRoom)
	mux.HandleFunc("/sendMessageOpenRoom", sendMessageOpenRoom)
	mux.HandleFunc("/sendMessageCloseRoom", sendMessageCloseRoom)
	mux.HandleFunc("/addToCloseRoom", addToCloseRoom)
	mux.HandleFunc("/removeMessage", removeMessage)
	mux.HandleFunc("/removeRoom", removeRoom)
	server := &http.Server{
		Addr: ":42069",
		Handler: mux,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	log.Println("Server running on http://localhost:42069")
	server.ListenAndServe()
}
