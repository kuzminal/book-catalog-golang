package frontend

import (
	"book-catalog/core"
	"book-catalog/storage"
	"encoding/json"
	"github.com/gin-gonic/gin"
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

/*func (f *RestFrontEnd) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}*/

func GetHandler(c *gin.Context) {
	key := c.Params.ByName("id")
	log.Println("Key : " + key)
	book, err := storage.GetBook(key)

	if err != nil {
		log.Fatal("Cannot get book from database")
	}

	if book != nil {
		renderJSON(c.Writer, book)
	} else {
		c.Writer.WriteHeader(http.StatusNotFound)
	}
}

func GetAllHandler(c *gin.Context) {
	book, err := storage.GetAll()
	if err != nil {
		log.Fatal("Cannot get books from database")
	}
	if book != nil {
		renderJSON(c.Writer, book)
	} else {
		c.Writer.WriteHeader(http.StatusNotFound)
	}
}

func DeleteHandler(c *gin.Context) {
	key := c.Params.ByName("id")
	DeleteBookChannel <- key
	c.Writer.WriteHeader(http.StatusAccepted)
}

func PostHandler(c *gin.Context) {
	// Enforce a JSON Content-Type.
	contentType := c.Request.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(c.Writer, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(c.Request.Body)
	dec.DisallowUnknownFields()
	var book core.Book
	if err := dec.Decode(&book); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}
	CreateBookChannel <- book
	c.Writer.WriteHeader(http.StatusAccepted)
}

func PutHandler(c *gin.Context) {
	// Enforce a JSON Content-Type.
	contentType := c.Request.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(c.Writer, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(c.Request.Body)
	dec.DisallowUnknownFields()
	var book core.Book
	if err := dec.Decode(&book); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}
	UpdateBookChannel <- book
	c.Writer.WriteHeader(http.StatusAccepted)
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
	r := gin.Default()

	r.GET(ApiVersion+ApiPrefix, GetAllHandler)
	r.POST(ApiVersion+ApiPrefix, PostHandler)
	r.GET(ApiVersion+ApiPrefix+"/:id", GetHandler)
	r.PUT(ApiVersion+ApiPrefix+"/:id", PutHandler)
	r.DELETE(ApiVersion+ApiPrefix+"/:id", DeleteHandler)

	return r.Run(AppPort)
}
