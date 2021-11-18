package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

const (
	ApiPrefix  = "/books"
	ApiVersion = "/v1"
	AppPort    = ":8081"
)

func asyncHandleChannels() {
	for {
		select {
		case book, ok := <-CreateBookChannel:
			if ok {
				err := SaveBook(&book)
				if err != nil {
					log.Fatal(err.Error())
				}
			}
		case book, ok := <-UpdateBookChannel:
			if ok {
				_, err := UpdateBook(book)
				if err != nil {
					log.Fatal(err.Error())
				}
			}
		case isbn, ok := <-DeleteBookChannel:
			if ok {
				err := DeleteBook(isbn)
				if err != nil {
					log.Fatal(err.Error())
				}
			}
		case <-time.After(4 * time.Second):
		}
	}
}

func main() {
	r := mux.NewRouter()

	go asyncHandleChannels()

	r.HandleFunc(ApiVersion+ApiPrefix, getAllHandler).Methods("GET")
	r.HandleFunc(ApiVersion+ApiPrefix, postHandler).Methods("POST")
	r.HandleFunc(ApiVersion+ApiPrefix+"/{key}", getHandler).Methods("GET")
	r.HandleFunc(ApiVersion+ApiPrefix+"/{key}", putHandler).Methods("PUT")
	r.HandleFunc(ApiVersion+ApiPrefix+"/{key}", deleteHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(AppPort, r))
}
