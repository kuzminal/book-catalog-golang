package main

import (
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
	clientOptions := options.Client().ApplyURI(mongoDbUri)
	client, err := mongo.Connect(ctx, clientOptions)
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

func SaveBook(book Book) error {
	//Perform InsertOne operation & validate against the error.
	_, err := collection.InsertOne(context.TODO(), book)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}

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
	return nil
}

func UpdateBook(book Book) (*Book, error) {
	//Perform InsertOne operation & validate against the error.
	filter := bson.D{primitive.E{Key: "isbn", Value: book.Isbn}}
	upsert := true
	after := options.After
	opt := options.FindOneAndReplaceOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
	result := collection.FindOneAndReplace(context.TODO(), filter, book, &opt)
	if result.Err() != nil {
		return nil, result.Err()
	}
	var t *Book
	fmt.Println(result.Decode(&t))
	return t, nil
}

func GetAll() ([]*Book, error) {
	var tasks []*Book
	filter := bson.D{primitive.E{Key: "", Value: nil}}
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return tasks, err
	}
	// Iterate through the cursor and decode each document one at a time
	for cur.Next(ctx) {
		var t Book
		err := cur.Decode(&t)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, &t)
	}
	if err := cur.Err(); err != nil {
		return tasks, err
	}
	// once exhausted, close the cursor
	cur.Close(ctx)
	if len(tasks) == 0 {
		return tasks, mongo.ErrNoDocuments
	}
	return tasks, nil
}

func GetBook(isbn string) (*Book, error) {
	var result *Book
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
	return result, nil
}
