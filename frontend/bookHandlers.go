package frontend

import (
	"book-catalog/core"
	"book-catalog/storage"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"mime"
	"net/http"
)

const (
	ApiPrefix  = "/books"
	ApiVersion = "/v1"
	AppPort    = ":8081"
)

type FrontEnd interface {
	Start() error
}

type RestFrontEnd struct {
}

var CreateBookChannel = make(chan core.Book, 16)
var UpdateBookChannel = make(chan core.Book, 16)
var DeleteBookChannel = make(chan string, 16)

func notAllowedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Allowed", http.StatusMethodNotAllowed)
}

func (f *RestFrontEnd) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Извлечь ключ из запроса
	key := vars["key"]
	book, err := storage.GetBook(key)

	if err != nil {
		log.Fatal("Cannot get book from database")
	}

	if book != nil {
		renderJSON(w, book)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func GetAllHandler(w http.ResponseWriter, r *http.Request) {
	book, err := storage.GetAll()
	if err != nil {
		log.Fatal("Cannot get books from database")
	}
	if book != nil {
		renderJSON(w, book)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Извлечь ключ из запроса
	key := vars["key"]
	DeleteBookChannel <- key
	w.WriteHeader(http.StatusAccepted)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
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
	var book core.Book
	if err := dec.Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	CreateBookChannel <- book
	w.WriteHeader(http.StatusAccepted)
}

func PutHandler(w http.ResponseWriter, r *http.Request) {
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
	var book core.Book
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

func (f *RestFrontEnd) Start() error {
	// Запомнить ссылку на основное приложение.

	r := mux.NewRouter()

	r.Use(f.loggingMiddleware)

	r.HandleFunc(ApiVersion+ApiPrefix, GetAllHandler).Methods("GET")
	r.HandleFunc(ApiVersion+ApiPrefix, PostHandler).Methods("POST")
	r.HandleFunc(ApiVersion+ApiPrefix+"/{key}", GetHandler).Methods("GET")
	r.HandleFunc(ApiVersion+ApiPrefix+"/{key}", PutHandler).Methods("PUT")
	r.HandleFunc(ApiVersion+ApiPrefix+"/{key}", DeleteHandler).Methods("DELETE")

	return http.ListenAndServe(AppPort, r)
}
