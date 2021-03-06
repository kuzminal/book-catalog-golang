package cache

import (
	"book-catalog/core"
	"context"
	"fmt"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"time"
)

var cacheBooks *cache.Cache
var keyCache = "books"
var ctx = context.TODO()
var RedisClient *redis.Client

func init() {
	//для кластера
	/*ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"localhost": ":6379",
		},
	})*/
	var redisHost = os.Getenv("REDIS_HOST")
	var redisPort = os.Getenv("REDIS_PORT")
	if redisHost == "" {
		redisHost = "localhost"
	}
	if redisPort == "" {
		redisPort = "6379"
	}
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	cacheBooks = cache.New(&cache.Options{
		Redis:      RedisClient,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})
}

func AddBookToCache(book *core.Book) {
	log.Printf("Add book to cache : %v", book)
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

func GetBookFromCache(isbn string) (*core.Book, error) {
	var book *core.Book
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
