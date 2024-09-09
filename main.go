package main

import (
	"github.com/nhannht/nhannht-have-xxx-stars/api"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/yeah", api.Handler)
	http.HandleFunc("/api/tada", api.TadaHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("you cannot even start a simple server, pathetic %v", err)
	}
}
