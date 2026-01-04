package main

import (
	"log"
	"net/http"
	"time"

	"github.com/FoldFunc/GoChat/server/app"
	"github.com/FoldFunc/GoChat/server/components"
	"github.com/FoldFunc/GoChat/server/db"
)
func main() {
	db.Init()
	mux := http.NewServeMux()

	mux.HandleFunc("/", components.Hello)

	mux.HandleFunc("/newUser", components.NewUser)

	mux.Handle("/newRoom", 
		app.AuthCookie(http.HandlerFunc(components.NewRoom)),
	)
	mux.Handle("/sendMessageOpenRoom", 
		app.AuthCookie(http.HandlerFunc(components.SendMessageOpenRoom)),
	)
	mux.Handle("/sendMessageCloseRoom", 
		app.AuthCookie(http.HandlerFunc(components.SendMessageCloseRoom)),
	)
	mux.Handle("/addToCloseRoom", 
		app.AuthCookie(http.HandlerFunc(components.AddToCloseRoom)),
	)
	mux.Handle("/addToOpenRoom", 
		app.AuthCookie(http.HandlerFunc(components.AddToOpenRoom)),
	)
	mux.Handle("/removeMessage", 
		app.AuthCookie(http.HandlerFunc(components.RemoveMessage)),
	)
	mux.Handle("/removeRoom", 
		app.AuthCookie(http.HandlerFunc(components.RemoveRoom)),
	)
	mux.Handle("/sendUserRequest", 
		app.AuthCookie(http.HandlerFunc(components.SendUserRequest)),
	)
	// Here
	mux.Handle("/viewUserRequests", 
		app.AuthCookie(http.HandlerFunc(components.ViewUserRequests)),
	)

	mux.Handle("/getNameById", 
		app.AuthCookie(http.HandlerFunc(components.GetNameById)),
	)

	mux.Handle("/queryUserRooms",
		app.AuthCookie(http.HandlerFunc(components.QueryUserRooms)),
	)
	mux.Handle("/queryUserChats",
		app.AuthCookie(http.HandlerFunc(components.QueryUserChats)),
	)
	mux.Handle("/queryUserChat",
		app.AuthCookie(http.HandlerFunc(components.QueryUserChat)),
	)

	server := &http.Server{
		Addr: ":42069",
		Handler: mux,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	log.Println("Server running on http://localhost:42069")
	server.ListenAndServe()
}
