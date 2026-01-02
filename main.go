/*
Also maybe make private and open rooms
So to some of them you can only get access if someone assgins you there.
*/
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"

)
var usedIds []int
type Message struct {
	Id     int    `json:"id"`
	UserId int    `json:"user_id"`
	Body   string `json:"body"`
}
type Type string
const (
	TypePublic  Type = "public"
	TypePrivate Type = "private"
)
type Room struct {
	Id       int       `json:"id"`
	UserId   int       `json:"user_id"`
	Name     string    `json:"name"`
	Messages []Message `json:"messages"`
	Type     Type      `json:"type"`
	Admins   []User    `json:"admins"`
	Users    []User    `json:"users"`
}
type User struct {
	Id    int     `json:"id"`
	Name  string  `json:"name"`
	Rooms []Room `json:"rooms"`
}
type GlobalMessages struct {
	Messages []Message `json:"messages"`
}
type GlobalUsers struct {
	Users []User `json:"users"`
}
type GlobalRooms struct {
	Rooms []Room `json:"rooms"`
}
type NewRoomReq struct {
	RoomId   int    `json:"room_id"`
	RoomName string `json:"room_name"`
	UserId   int    `json:"user_id"`
	RoomType bool   `json:"room_type"`
}
type NewUserReq struct {
	UserName string `json:"user_name"`
}
type SendMessreq struct {
	RoomId int    `json:"room_id"`
	UserId int    `json:"user_id"`
	Body   string `json:"body"`
}
type RemovemesReq struct {
	RoomId int    `json:"room_id"`
	UserId int    `json:"user_id"`
	MessId int    `json:"mess_id"`
}
type RemoveRoomReq struct {
	RoomId int    `json:"room_id"`
	UserId int    `json:"user_id"`
}
type AddToCloseRoomReq struct {
	RoomId  int `json:"room_id"`
	AdminId int `json:"admin_id"`
	UserId  int `json:"user_id"`
}
type AccesRoomReq struct {
	RoomId int `json:"room_id"`
	UserId int `json:"user_id"`
}
func hello(w http.ResponseWriter, r *http.Request) {
	log.Println("/ handler called")
}
func generateId() (int) {
	var cont bool
	var newId int
	for {
		cont = false
		newId = rand.IntN(1 << 32)
		for _, n := range usedIds {
			if n == newId {
				cont = true
			}
		}
		if !cont {
			break
		}
	}
	usedIds = append(usedIds, newId)
	return newId
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
func IsAdmin(user_id, room_id int) (bool) {
	for _, r := range R.Rooms {
		if r.Id == room_id {
			for _, a := range r.Admins {
				if a.Id == user_id {
					return true
				}
			}
		}
	}
	return false
}
func messageExsists(room_id, mess_id, user_id int) (bool) {
	for _, r := range R.Rooms {
		if r.Id == room_id {
			for _, m := range r.Messages {
				if m.Id == mess_id && m.UserId == user_id {
					return true
				}
			}
		}
	}
	return false
}
func roomPublic(id int) (bool) {
	for _, u := range R.Rooms {
		if u.Id == id && u.Type == TypePublic {
			return true
		}
	}
	return false
}
func roomExsists(id int) (bool) {
	for _, u := range R.Rooms {
		if u.Id == id {
			return true
		}
	}
	return false
}
func roomExsistsToDelete(id, user_id int) (bool) {
	for _, u := range R.Rooms {
		if u.Id == id && u.UserId == user_id {
			return true
		}
	}
	return false
}
func userExsists(id int) (bool) {
	for _, u := range U.Users {
		if u.Id == id {
			return true
		}
	}
	return false
}
func UserInRoom(user_id, room_id int) (bool) {
	for _, u := range U.Users {
		if u.Id == user_id {
			for _, r := range u.Rooms {
				if r.Id == room_id {
					return true
				}
			}
		}
	}
	return false
}
func getUserById(user_id int) (User, error) {
	for _, u := range U.Users {
		if u.Id == user_id {
			return u, nil
		}
	}
	return User{}, fmt.Errorf("Invalid user id")
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
	for _, r := range R.Rooms {
		fmt.Printf("Room id: %d;Room name: %s;Room public: %v\n", newId, r.Name, r.Type)
		if req.RoomName == r.Name {
			http.Error(w, "Room name already exsists", http.StatusConflict)
			return
		}
	}
	R.Rooms = append(R.Rooms, NewRoom)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"id": newId,
	})
}
// Depricated
func getRoom(w http.ResponseWriter, r *http.Request) {
	log.Println("/getRoom called")
	if r.Method != http.MethodGet {
		http.Error(w, "Only get requests allowed", http.StatusMethodNotAllowed)
		return
	}
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	newId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	if !userExsists(newId) {
		http.Error(w, "User validation failed", http.StatusForbidden)
		return
	}
	var returnRoom Room
	for _, r := range R.Rooms {
		if r.Name == name{
			returnRoom = r
		}
	}
	idInt := strconv.Itoa(returnRoom.Id)
	fmt.Printf("Room id requested: %d\n", returnRoom.Id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id": idInt,
		"type": string(returnRoom.Type),
	})
}
func sendMessageCloseRoom(w http.ResponseWriter, r *http.Request) {
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
	for _, r := range R.Rooms {
		if r.Id == req.RoomId {
			r.Messages = append(r.Messages, message)
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
	for _, r := range R.Rooms {
		if r.Id == req.RoomId {
			r.Messages = append(r.Messages, message)
		}
	}
	idInt := strconv.Itoa(id)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "message created",
		"id": idInt,
	})
}
func removeMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	var req RemovemesReq
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
		http.Error(w, "No such room", http.StatusBadRequest)
		return
	}
	if !messageExsists(req.RoomId, req.MessId, req.UserId) {
		http.Error(w, "No such message", http.StatusBadRequest)
		return
	}
	for _, r := range R.Rooms {
		if r.Id == req.RoomId {
			for i, m := range r.Messages {
				if m.Id == req.MessId {
					r.Messages = append(r.Messages[:i], r.Messages[i+1:]...)
				}
			}
		}
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "message removed",
	})
}
func removeRoom(w http.ResponseWriter, r *http.Request) {
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
	for i, r := range R.Rooms {
		if r.Id == req.RoomId {
			R.Rooms = append(R.Rooms[:i], R.Rooms[i+1:]...)
		}
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Room deleted",
	})
}
func addToCloseRoom(w http.ResponseWriter, r *http.Request) {
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
	for _, r := range R.Rooms {
		if r.Id == req.RoomId {
			r.Users = append(r.Users, currentUser)
		}
	}
}
func accesRoom(w http.ResponseWriter, r *http.Request) {
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
	for _, r := range R.Rooms {
		if r.Id == req.RoomId {
			r.Users = append(r.Users, currentUser)
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
	mux.HandleFunc("/getRoom", getRoom)
	mux.HandleFunc("/accesRoom", accesRoom)
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
