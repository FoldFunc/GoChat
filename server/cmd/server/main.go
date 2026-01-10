package main

import (
	"log"
	"os"

	"github.com/FoldFunc/GoChat/server/db"
	apphttp "github.com/FoldFunc/GoChat/server/internal/http"
)
func main() {
	addr := os.Getenv("PORT")
	if addr == "" {
		addr = "8080"
	}
	addr = ":" + addr
	handler := apphttp.Routes()
	server := apphttp.New(addr, handler, "", "")
	db.Init()
	log.Printf("listening on port: %s (HTTP)", addr)
	log.Fatal(server.Start())
}
