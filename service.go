package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	ApiPrefix  = "/books"
	ApiVersion = "/v1"
	AppPort    = ":8081"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc(ApiVersion+ApiPrefix, getAllHandler).Methods("GET")
	r.HandleFunc(ApiVersion+ApiPrefix+"/{key}", getHandler).Methods("GET")
	r.HandleFunc(ApiVersion+ApiPrefix, postHandler).Methods("POST")
	r.HandleFunc(ApiVersion+ApiPrefix+"/{key}", putHandler).Methods("PUT")

	log.Fatal(http.ListenAndServe(AppPort, r))
}
