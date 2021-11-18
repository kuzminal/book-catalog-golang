package main

import (
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var cacheBooks *cache.Cache
var keyCache = "books"

func init() {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"localhost": ":6379",
		},
	})

	cacheBooks = cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})
}

func AddBookToCache(book *Book) {
	if err := cacheBooks.Set(&cache.Item{
		Ctx:   ctx,
		Key:   keyCache + ":" + book.Isbn,
		Value: book,
		TTL:   time.Hour,
	}); err != nil {
		panic(err)
	}
	log.Printf("Cached book : %v", book)
}

func GetBookFromCache(isbn string) (*Book, error) {
	var book *Book
	if err := cacheBooks.Get(ctx, keyCache+":"+isbn, &book); err != nil {
		log.Printf("Cannot find book with isbn : %s", isbn)
		return nil, err
	}
	log.Printf("Got from cache book : %v", book)
	return book, nil
}

func DeleteBookFromCache(isbn string) {
	if err := cacheBooks.Delete(ctx, keyCache+":"+isbn); err != nil {
		panic(err)
	}
	log.Printf("Deleted from cache book : %v", isbn)
}
