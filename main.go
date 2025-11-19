package main

import (
	"net/http"

	"github.com/21yyri/search-engine/handlers"
)

func main() {
	server := http.NewServeMux()

	server.HandleFunc("/add", handlers.RegisterUrl)
	server.HandleFunc("/get", handlers.SearchUrls)

	http.ListenAndServe(":8080", server)
}
