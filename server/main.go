package main

import (
	"log"
	"net/http"
	"time"
	"github.com/FoldFunc/GoChat/server/components"
)
func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", components.Hello)

	mux.HandleFunc("/newUser", components.NewUser)

	mux.HandleFunc("/newRoom", components.NewRoom)

	mux.HandleFunc("/sendMessageOpenRoom", components.SendMessageOpenRoom)
	mux.HandleFunc("/sendMessageCloseRoom", components.SendMessageCloseRoom)

	mux.HandleFunc("/addToCloseRoom", components.AddToCloseRoom)
	mux.HandleFunc("/addToOpenRoom",components.AddToOpenRoom)

	mux.HandleFunc("/removeMessage", components.RemoveMessage)
	mux.HandleFunc("/removeRoom", components.RemoveRoom)

	mux.HandleFunc("/sendUserRequest", components.SendUserRequest)
	mux.HandleFunc("/viewUserRequests", components.ViewUserRequests)

	mux.HandleFunc("/getNameById", components.GetNameById)
	server := &http.Server{
		Addr: ":42069",
		Handler: mux,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	log.Println("Server running on http://localhost:42069")
	server.ListenAndServe()
}
