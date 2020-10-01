package helpers

import (
	"encoding/json"
	"net/http"
)

func ResponseJSON(rw http.ResponseWriter, body interface{}, code int) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)

	if nil != body {
		json.NewEncoder(rw).Encode(body)
	}
}
