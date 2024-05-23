package server

import (
	"log"
	"net/http"
)

func HttpError(w http.ResponseWriter, msg string, code int, err error) {
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(code)
	if _, err := w.Write([]byte(msg)); err != nil {
		log.Fatalf("failed to write http response: %v\n", err)
	}
}
