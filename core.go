package main

type Book struct {
	//Id             string  `bson:"_id,omitempty"`
	Isbn           string  `bson:"isbn,omitempty"`
	Title          string  `bson:"title,omitempty"`
	Author         string  `bson:"author,omitempty"`
	PublishingYear int     `bson:"publishingYear,omitempty"`
	Price          float64 `bson:"price,omitempty"`
	Quantity       int     `bson:"quantity,omitempty"`
	Version        int     `bson:"version,omitempty"`
}
