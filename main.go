package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/21yyri/search-engine/models"
	"github.com/21yyri/search-engine/utils"
)

func main() {
	server := http.NewServeMux()

	server.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		var url models.Url
		json.NewDecoder(r.Body).Decode(&url)

		err := utils.TokenizeResponse(url.Url)
		if err != nil {
			fmt.Fprintf(w, "%s", err.Error())
			return
		}

		fmt.Fprint(w, http.StatusCreated)
	})

	server.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		var query models.Query
		json.NewDecoder(r.Body).Decode(&query)

		response, _ := json.Marshal(
			models.QueryResponse{
				Results: utils.FetchUrls(query.Query),
			},
		)

		fmt.Fprint(w, string(response))
	})

	http.ListenAndServe(":8080", server)
}
