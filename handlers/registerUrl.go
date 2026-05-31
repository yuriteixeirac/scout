package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/21yyri/search-engine/models"
	"github.com/21yyri/search-engine/utils"
)

func RegisterUrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.Urls
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}

	err = utils.TokenizeResponse(req.Url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)

	fmt.Fprint(w)
}
