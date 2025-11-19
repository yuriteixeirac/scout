package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/21yyri/search-engine/models"
	"github.com/21yyri/search-engine/utils"
)

func RegisterUrl(w http.ResponseWriter, r *http.Request) {
	var urls models.Urls
	json.NewDecoder(r.Body).Decode(&urls)

	for _, url := range urls.Url {
		err := utils.TokenizeResponse(url)
		if err != nil {
			fmt.Fprintf(w, "%s", err.Error())
			return
		}
	}

	fmt.Fprint(w, http.StatusCreated)
}
