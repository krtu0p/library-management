package dao

import "time"

type Borrowing struct {
	ID				uint		`gorm:"primarykey"`
	BorrowDate 		*time.Time
	ReturnDate 		*time.Time
	BookID 			uint 		`gorm:"not null;"`
	BorrowedBook 	*Book		`gorm:"foreignKey:BookID;"`
	PersonID 		uint		`gorm:"not null;"`
	BorrowerPerson 	*Person		`gorm:"foreignKey:PersonID;"`
}

func (Borrowing) TableName() string {
	return "borrowings"
}