package components

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/FoldFunc/GoChat/server/app"
	"github.com/FoldFunc/GoChat/server/db"
)

func QueryUserRooms(w http.ResponseWriter, r *http.Request) {
	log.Println("/queryUserRooms called")
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	userId := r.Context().Value("userID").(int)
	rooms, err := db.QuerUserRoomsDB(userId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Println(rooms)
	w.Header().Set("Content-Type", "application/json")	
	json.NewEncoder(w).Encode(rooms)
}
func QueryUserChats(w http.ResponseWriter, r *http.Request) {
	log.Println("/queryUserChats called")
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	userId := r.Context().Value("userID").(int)
	chats, err := db.QueryUserChatsDB(userId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Println(chats)
	w.Header().Set("Content-Type", "application/json")	
	json.NewEncoder(w).Encode(chats)
}
func QueryUserChat(w http.ResponseWriter, r *http.Request) {
	log.Println("/queryUserChat called")
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
	chat, err := db.QuerySpecificUserChatDB(userId, req.ChatWithName)
	log.Println(chat)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chat)
}

func QueryUserRoom(w http.ResponseWriter, r *http.Request) {
	log.Println("/queryUserRoom called")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	var req app.QueryUserRoomReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("HERE")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	room, err := db.QueryPublicRoomByNameDB(req.RoomName)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Println(room)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room)
}
func QueryMessageFromRoom(w http.ResponseWriter, r *http.Request) {
	log.Println("/queryMessageFromRoom called")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	var req app.QueryMessageFromRoomReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("JSON ERROR: ", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	userId := r.Context().Value("userID").(int)
	exsists, err := db.RoomExistsDB(req.RoomId)
	if err != nil {
		log.Println("RoomExistsDB error: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !exsists {
		http.Error(w, "Room does not exsist", http.StatusBadRequest)
		return
	}
	messageId, err := db.GetMessageIDByNameDB(req.RoomId, req.MessageBody, userId)
	if err != nil {
		log.Println("GetMessageIDByNameDB error: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	messageExsists, err := db.MessageExistsDB(req.RoomId, messageId, userId)
	if err != nil {
		log.Println("MessageExistsDB error: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !messageExsists {
		http.Error(w, "Message does not exsist", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"message_id": messageId,
	})
}
