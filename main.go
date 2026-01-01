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
type GlobalMessages struct {
	Messages []Message `json:"messages"`
}
type NewMessageReq struct {
	Message string `json:"message"`
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
var M GlobalMessages 

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	mux.HandleFunc("/newMessage", newMessage)
	mux.HandleFunc("/getMessage", getMessage)
	mux.HandleFunc("/getMessages", getMessages)
	server := &http.Server{
		Addr: ":42069",
		Handler: mux,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	log.Println("Server running on http://localhost:42069")
	server.ListenAndServe()
}
