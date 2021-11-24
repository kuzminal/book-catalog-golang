package main

import (
	"book-catalog/frontend"
	"book-catalog/storage"
	"log"
	"time"
)

func asyncHandleChannels() {
	for {
		select {
		case book, ok := <-frontend.CreateBookChannel:
			if ok {
				err := storage.SaveBook(&book)
				if err != nil {
					log.Fatal(err.Error())
				}
			}
		case book, ok := <-frontend.UpdateBookChannel:
			if ok {
				_, err := storage.UpdateBook(book)
				if err != nil {
					log.Fatal(err.Error())
				}
			}
		case isbn, ok := <-frontend.DeleteBookChannel:
			if ok {
				err := storage.DeleteBook(isbn)
				if err != nil {
					log.Fatal(err.Error())
				}
			}
		case <-time.After(4 * time.Second):
		}
	}
}

func main() {

	go asyncHandleChannels()
	fe := frontend.RestFrontEnd{}
	fe.Start()

}
