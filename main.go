package main

import (
	"fmt"
	"net/http"

	"github.com/21yyri/search-engine/handlers"
)

func main() {
	server := http.NewServeMux()

	server.HandleFunc("POST /add", handlers.RegisterUrl)
	server.HandleFunc("POST /query", handlers.SearchUrls)

	err := http.ListenAndServe(":8080", server)
	if err != nil {
		fmt.Println(err.Error())
	}
}
