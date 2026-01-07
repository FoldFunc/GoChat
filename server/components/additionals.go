package components

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/FoldFunc/GoChat/server/app"
	"github.com/FoldFunc/GoChat/server/db"
)
func Hello(w http.ResponseWriter, r *http.Request) {
	log.Println("/ handler called")
}
func GetNameById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	var req app.GetNameByIdReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	userId := r.Context().Value("userID").(int)
	userIs, err := db.UserExists(userId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !userIs {
		http.Error(w, "Cookie invalid", http.StatusForbidden)
		return
	}
	name, err := db.GetNameByIdDB(userId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Println("Name: ", name)
	json.NewEncoder(w).Encode(map[string]string{
		"name": name,
	})
}
func GetIdByName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	var req app.GetIdByNameReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	userId := r.Context().Value("userID").(int)
	userIs, err := db.UserExists(userId)
	if err != nil {
		log.Println("CALLED IT 0")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !userIs {
		http.Error(w, "Cookie invalid", http.StatusForbidden)
		return
	}
	id, err := db.GetUserIdByNameDB(req.SearchName)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Println("Id: ", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id": strconv.Itoa(id),
	})
}
