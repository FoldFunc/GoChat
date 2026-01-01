package main

import (
	"encoding/json"
	"fmt"
	"log"
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
type GlobalMessages struct {
	Messages []Message `json:"messages"`
}
type GlobalRooms struct {
	Rooms []Room `json:"rooms"`
}
type NewMessageReq struct {
	Message string `json:"message"`
}
type NewRoomReq struct {
	RoomId int `json:"room_id"`
	RoomName string `json:"room_name"`
}
func hello(w http.ResponseWriter, r *http.Request) {
	log.Println("/ handler called")
}
func strint(value string) (int, error) {
	i, err := strconv.Atoi(value)
	if err != nil {
		return -1, err
	}
	return i, nil
}
func getMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only get requests allowed", http.StatusMethodNotAllowed)
		return
	}
	meesages := M.Messages
	chat := GlobalMessages{
		Messages: meesages,
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(chat)
	if err != nil {
		http.Error(w, "Failed to encode messages", http.StatusInternalServerError)
		return
	}
	
}
func getMessage(w http.ResponseWriter, r *http.Request) {
	log.Println("/getMessage called")
	if r.Method != http.MethodGet {
		http.Error(w, "Only get requests allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	idNew, err := strint(id)
	if err != nil {
		http.Error(w, "Invalid number", http.StatusBadRequest)
		return
	}
	var returnMessage Message
	for _, m := range M.Messages {
		if m.Id == idNew {
			returnMessage = m
		}
	}
	fmt.Printf("Message requested: %s", returnMessage.Body)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(returnMessage)
	if err != nil {
		http.Error(w, "Failed to encode messages", http.StatusInternalServerError)
		return
	}
}
func newMessage(w http.ResponseWriter, r *http.Request) {
	log.Println("/newMessage called")
	if r.Method != http.MethodPost {
		http.Error(w, "Only post requests allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var req NewMessageReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Recived message: %s\n", req.Message)
	var newId int
	if len(M.Messages) == 0 {
		newId = 0
	} else {
		lastMessage := M.Messages[len(M.Messages)-1]
		newId = lastMessage.Id + 1
	}
	NewMessage := Message{
		Id: newId,
		Body: req.Message,
	}
	M.Messages = append(M.Messages, NewMessage)
}
func createRoom(w http.ResponseWriter, r *http.Request) {
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

	fmt.Fprintf(w, "Recived creating room request of name: %s\n", req.RoomName)
	var newId int
	if len(R.Rooms) == 0 {
		newId = 0
	} else {
		lastMessage := R.Rooms[len(R.Rooms)-1]
		newId = lastMessage.Id + 1
	}
	NewRoom := Room{
		Id: newId,
		Name: req.RoomName,
	}
	R.Rooms = append(R.Rooms, NewRoom)
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
	var returnRoom Room
	for _, r := range R.Rooms {
		if r.Name == name{
			returnRoom = r
		}
	}
	fmt.Printf("Room id requested: %d", returnRoom.Id)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(returnRoom.Id)
	if err != nil {
		http.Error(w, "Failed to encode messages", http.StatusInternalServerError)
		return
	}
}
var R GlobalRooms
var M GlobalMessages 

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	mux.HandleFunc("/newMessage", newMessage)
	mux.HandleFunc("/getMessage", getMessage)
	mux.HandleFunc("/getMessages", getMessages)
	mux.HandleFunc("/newRoom", createRoom)
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
