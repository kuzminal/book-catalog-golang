package core

type Book struct {
	//Id             string  `bson:"_id,omitempty"`
	Isbn           string  `bson:"isbn,omitempty" json:"isbn,omitempty"`
	Title          string  `bson:"title,omitempty" json:"title,omitempty"`
	Author         string  `bson:"author,omitempty" json:"author,omitempty"`
	PublishingYear int     `bson:"publishingYear,omitempty" json:"publishingYear,omitempty"`
	Price          float64 `bson:"price,omitempty" json:"price,omitempty"`
	Quantity       int     `bson:"quantity,omitempty" json:"quantity,omitempty"`
	Version        int     `bson:"version,omitempty" json:"version,omitempty"`
}
