package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/21yyri/search-engine/models"
	"github.com/21yyri/search-engine/utils"
)

func SearchUrls(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var query models.Query
	err := json.NewDecoder(r.Body).Decode(&query)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w)
		return
	}

	queryResults, err := utils.FetchUrls(query.Query)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}

	response, err := json.Marshal(
		models.QueryResponse{
			Results: queryResults,
		},
	)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(response))
}
