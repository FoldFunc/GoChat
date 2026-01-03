package components

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/FoldFunc/GoChat/server/app"
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
	if !app.UserExsists(userId) {
		http.Error(w, "User verification failed", http.StatusForbidden)
		return
	}
	if !app.UserExsists(req.SearchId) {
		http.Error(w, "User verification failed", http.StatusForbidden)
		return
	}
	var name string
	for _, u := range app.U.Users {
		if u.Id == req.SearchId {
			name = u.Name
		}
	}
	json.NewEncoder(w).Encode(map[string]string{
		"name": name,
	})
}
