package components

import (
	"encoding/json"
	"net/http"

	"github.com/FoldFunc/GoChat/server/app"
)

func QueryUserRooms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	userId := r.Context().Value("userID").(int)
	var rooms []*app.Room
	for _, u := range app.U.Users {
		if u.Id == userId {
			rooms = u.Rooms
		}
	}
	w.Header().Set("Content-Type", "application/json")	
	json.NewEncoder(w).Encode(rooms)
}
func QueryUserChats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	userId := r.Context().Value("userID").(int)
	var chats []*app.Chat
	for _, u := range app.U.Users {
		if u.Id == userId {
			chats = u.Chats
		}
	}
	w.Header().Set("Content-Type", "application/json")	
	json.NewEncoder(w).Encode(chats)
}
func QueryUserChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	userId := r.Context().Value("userID").(int)
	var req app.QueryUserChatReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	user, err := app.GetUserById(userId)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}
	var chat app.Chat
	for _, c := range user.Chats {
		if c.User2.Name == req.ChatWithName {
			chat = *c
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chat)
}
func QueryUserRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	userId := r.Context().Value("userID").(int)
	var req app.QueryUserRoomReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	user, err := app.GetUserById(userId)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}
	var room app.Room
	for _, r := range user.Rooms {
		if r.Name == req.RoomName {
			room = *r
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room)
}
