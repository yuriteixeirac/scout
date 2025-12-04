package main

import (
	"net/http"

	"github.com/21yyri/search-engine/handlers"
)

func main() {
	server := http.NewServeMux()

	server.HandleFunc("POST /add", handlers.RegisterUrl)
	server.HandleFunc("GET /get", handlers.SearchUrls)

	http.ListenAndServe(":8080", server)
}
