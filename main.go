package main

import (
	"github.com/nhannht/nhannht-have-xxx-stars/app/api"
	"net/http"
)

func main() {
	http.HandleFunc("/api/yeah", api.Handler)
	http.ListenAndServe(":8080", nil)
}
