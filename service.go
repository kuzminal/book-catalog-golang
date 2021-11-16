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

func asyncWrite() {
	book, ok := <-CreateBookChannel
	if ok {
		err := SaveBook(book)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func asyncUpdate() {
	book, ok := <-UpdateBookChannel
	if ok {
		_, err := UpdateBook(book)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func asyncDelete() {
	isbn, ok := <-DeleteBookChannel
	if ok {
		err := DeleteBook(isbn)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func main() {
	r := mux.NewRouter()

	go asyncWrite()
	go asyncUpdate()
	go asyncDelete()

	r.HandleFunc(ApiVersion+ApiPrefix, getAllHandler).Methods("GET")
	r.HandleFunc(ApiVersion+ApiPrefix, postHandler).Methods("POST")
	r.HandleFunc(ApiVersion+ApiPrefix+"/{key}", getHandler).Methods("GET")
	r.HandleFunc(ApiVersion+ApiPrefix+"/{key}", putHandler).Methods("PUT")
	r.HandleFunc(ApiVersion+ApiPrefix+"/{key}", deleteHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(AppPort, r))
}
