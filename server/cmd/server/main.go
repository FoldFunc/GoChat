package main

import (
	"log"
	"os"

	"github.com/FoldFunc/GoChat/server/db"
	apphttp "github.com/FoldFunc/GoChat/server/internal/http"
)
func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	handler := apphttp.Routes()
	server := apphttp.New(addr, handler)
	db.Init()
	log.Printf("listening on port: %s", addr)
	log.Fatal(server.Start())
}
