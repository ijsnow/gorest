package gorest

import (
	"encoding/json"
	"io"
	"net/http"
)

// WriteFunc is the type for the function used to write data in the way we want
type WriteFunc func(http.ResponseWriter, int, interface{})

// JSON writes data in JSON format
func JSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")

	bytes, err := json.Marshal(data)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Internal server error")
	} else {
		w.WriteHeader(code)
		io.WriteString(w, string(bytes))
	}
}
