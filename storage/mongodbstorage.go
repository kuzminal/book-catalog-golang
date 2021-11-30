package storage

import (
	"book-catalog/cache"
	"book-catalog/core"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var collection *mongo.Collection
var ctx = context.TODO()
var mongoDbUri = os.Getenv("MONGODB_URI")
var mongoDbDatabase = os.Getenv("MONGODB_DATABASE")
var mongoDbCollection = os.Getenv("MONGODB_COLLECTION")

func init() {
	if mongoDbUri == "" {
		mongoDbUri = "mongodb://localhost:27017/"
	}
	opt := options.ClientOptions{}
	opt.ApplyURI(mongoDbUri)
	timeOut := 30 * time.Second
	opt.ConnectTimeout = &timeOut
	opt.ServerSelectionTimeout = &timeOut
	//clientOptions := options.Client().ApplyURI(mongoDbUri)
	cont, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	client, err := mongo.Connect(cont, &opt)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	if mongoDbDatabase == "" {
		mongoDbDatabase = "book-catalog"
	}
	if mongoDbCollection == "" {
		mongoDbCollection = "books"
	}

	collection = client.Database("book-catalog").Collection("books")
}

/*
func SaveBook(book *core.Book) error {
	//Perform InsertOne operation & validate against the error.
	upsert := true
	_, err := collection.UpdateOne(context.TODO(), book, options.UpdateOptions{
		Upsert: &upsert,
	})
	if err != nil {
		log.Println("save book")
		log.Println(err.Error())
		return err
	}
	//Return success without any error.
	go cache.AddBookToCache(book)
	return nil
}*/

func DeleteBook(isbn string) error {
	filter := bson.D{primitive.E{Key: "isbn", Value: isbn}}
	duration := time.Second
	opt := options.FindOneAndDeleteOptions{
		MaxTime: &duration,
	}
	result := collection.FindOneAndDelete(context.TODO(), filter, &opt)
	if result.Err() != nil {
		return result.Err()
	}
	//Return success without any error.
	go cache.DeleteBookFromCache(isbn)
	return nil
}

func UpdateBook(book *core.Book) (*core.Book, error) {
	//Perform InsertOne operation & validate against the error.
	filter := bson.D{primitive.E{Key: "isbn", Value: book.Isbn}}
	upsert := true
	after := options.After
	opt := options.FindOneAndReplaceOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
	var bookDecoded *core.Book
	result := collection.FindOneAndReplace(context.TODO(), filter, book, &opt)
	if result.Err() != nil {
		return nil, result.Err()
	}
	fmt.Println(result.Decode(&bookDecoded))
	go cache.AddBookToCache(bookDecoded)
	return bookDecoded, nil
}

func GetAll() ([]*core.Book, error) {
	var books []*core.Book
	filter := bson.D{primitive.E{Key: "", Value: nil}}
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return books, err
	}
	// Iterate through the cursor and decode each document one at a time
	for cur.Next(ctx) {
		var book core.Book
		err := cur.Decode(&book)
		if err != nil {
			return books, err
		}
		books = append(books, &book)
	}
	if err := cur.Err(); err != nil {
		return books, err
	}
	// once exhausted, close the cursor
	cur.Close(ctx)
	if len(books) == 0 {
		return books, mongo.ErrNoDocuments
	}
	return books, nil
}

func GetBook(isbn string) (*core.Book, error) {
	book, errCache := cache.GetBookFromCache(isbn)
	if errCache != nil {
		var result *core.Book
		filter := bson.D{primitive.E{Key: "isbn", Value: isbn}}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := collection.FindOne(ctx, filter).Decode(&result)
		if err == mongo.ErrNoDocuments {
			// Do something when no record was found
			fmt.Println("record does not exist")
		} else if err != nil {
			log.Fatal(err)
			return result, err
		}
		go cache.AddBookToCache(result)
		return result, nil
	} else {
		return book, nil
	}
}
