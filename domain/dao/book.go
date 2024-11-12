package dao

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title 			string 		`gorm:"size:56;not null;"`
	Subtitle 		string 		`gorm:"size:64;not null;"`
	AuthorID 		uint 		`gorm:"not null;"`
	BookAuthor  	*Author		`gorm:"foreignKey:AuthorID;"`
	PublisherID 	uint 		`gorm:"not null;"`
	BookPublisher   *Publisher	`gorm:"foreignKey:PublisherID;"`
}

func (Book) TableName() string {
	return "books"
}