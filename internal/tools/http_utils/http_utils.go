package http_utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func SetJSONResponse(w http.ResponseWriter, body interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")

	result, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte("{\"error\": \"can't marshal body\"}")); err != nil {
			log.Fatal(err)
		}
		return
	}
	w.WriteHeader(statusCode)
	if _, err := w.Write(result); err != nil {
		log.Fatal(err)
	}
}
