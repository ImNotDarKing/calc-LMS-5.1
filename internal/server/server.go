package server

import (
	v1 "calc-LMS/internal/api/v1"
	"log"
	"net/http"
)


func StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/calculate", v1.Calculate)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
