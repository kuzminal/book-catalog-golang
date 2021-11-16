package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"mime"
	"net/http"
)

var CreateBookChannel = make(chan Book, 16)
var UpdateBookChannel = make(chan Book, 16)
var DeleteBookChannel = make(chan string, 16)

func notAllowedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Allowed", http.StatusMethodNotAllowed)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Извлечь ключ из запроса
	key := vars["key"]
	book, err := GetBook(key)

	if err != nil {
		log.Fatal("Cannot get book from database")
	}

	if book != nil {
		renderJSON(w, book)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func getAllHandler(w http.ResponseWriter, r *http.Request) {
	book, err := GetAll()
	if err != nil {
		log.Fatal("Cannot get books from database")
	}
	if book != nil {
		renderJSON(w, book)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Извлечь ключ из запроса
	key := vars["key"]
	DeleteBookChannel <- key
	w.WriteHeader(http.StatusAccepted)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	// Enforce a JSON Content-Type.
	contentType := r.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var book Book
	if err := dec.Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	CreateBookChannel <- book
	w.WriteHeader(http.StatusAccepted)
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	// Enforce a JSON Content-Type.
	contentType := r.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var book Book
	if err := dec.Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	/*result, err := UpdateBook(book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	renderJSON(w, result)*/
	UpdateBookChannel <- book
	w.WriteHeader(http.StatusAccepted)
}

// renderJSON renders 'v' as JSON and writes it as a response into w.
func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
