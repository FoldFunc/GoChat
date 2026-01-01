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
type Message struct {
	Id int `json:"id"`
	Body string `json:"body"`
}
type Room struct {
	Id int `json:"id"`
	Name string `json:"name"`
}
type User struct {
	Id int `json:"id"`
	Name string `json:"name"`
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
	RoomId int `json:"room_id"`
	RoomName string `json:"room_name"`
	UserId int `json:"user_id"`
}
type NewUserReq struct {
	UserName string `json:"user_name"`
}
func hello(w http.ResponseWriter, r *http.Request) {
	log.Println("/ handler called")
}
func generateId() (int) {
	return rand.IntN(1 << 32)
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
func userExsists(id int) (bool) {
	for _, u := range U.Users {
		if u.Id == id {
			return true
		}
	}
	return false
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
	newId := generateId()
	NewRoom := Room{
		Id: newId,
		Name: req.RoomName,
	}
	for _, r := range R.Rooms {
		fmt.Printf("Room id: %d;Room name: %s\n", newId, r.Name)
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
	fmt.Printf("Room id requested: %d\n", returnRoom.Id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"id": returnRoom.Id,
	})
}
var R GlobalRooms
var M GlobalMessages 
var U GlobalUsers
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	mux.HandleFunc("/newUser", newUser)
	mux.HandleFunc("/newRoom", newRoom)
	mux.HandleFunc("/getRoom", getRoom)
	server := &http.Server{
		Addr: ":42069",
		Handler: mux,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	log.Println("Server running on http://localhost:42069")
	server.ListenAndServe()
}
