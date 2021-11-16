package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"log"
	"mime"
	"net/http"
	"strings"
)

func notAllowedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Allowed", http.StatusMethodNotAllowed)
}

func checkToken(r *http.Request) error {
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.ReplaceAll(authHeader, "Bearer ", "")
	log.Println(tokenString)
	//tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU"

	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte("MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAhAZou/UbaPU5O7uSjvS4CmtN6Dk9bY/MlwwvtQ5IjsgduPiRWz4gQtpp6LiG9yvkGnKoQOXYB63N/7sNoqUeMB/AIICY4blFDX+/mWs4n/uGa3APOIItkqLz4E4Dix4UmPxSjd5qg73GjP4yPTH9VQq5kfzcw3ohHGk9RrpeUEE3wmB93uOunNOSLnDHnY/4Ssy8/uKY6Ua6T3dDWLir7EApyPlhlfHbgrWd6vsMIuDBiUwVYCvqtcBFbD1gSWpk0j84CFhqzryrCzECklNCjQBrCH98gxJBrM4zcKWeB8uVKzRa6Qc5tTLasBE9nwzh4aKfDQaGkkzKXhseWXmkAwIDAQAB"))
	if err != nil {
		return fmt.Errorf("validate: parse key: %w", err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return key, nil
	})

	_, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		fmt.Println(err)
	}
	return nil
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	checkToken(r)
	vars := mux.Vars(r) // Извлечь ключ из запроса
	key := vars["key"]
	book, err := GetBook(key)
	if err != nil {
		log.Fatal("Cannot get book from database")
	}
	renderJSON(w, book)
}

func getAllHandler(w http.ResponseWriter, r *http.Request) {
	book, err := GetAll()
	if err != nil {
		log.Fatal("Cannot get books from database")
	}
	renderJSON(w, book)
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

	err = SaveBook(book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

	result, err := UpdateBook(book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	renderJSON(w, result)
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
