package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/21yyri/search-engine/models"
	"github.com/21yyri/search-engine/utils"
)

func SearchUrls(w http.ResponseWriter, r *http.Request) {
	var query models.Query
	json.NewDecoder(r.Body).Decode(&query)

	response, _ := json.Marshal(
		models.QueryResponse{
			Results: utils.FetchUrls(query.Query),
		},
	)

	fmt.Fprint(w, string(response))
}
