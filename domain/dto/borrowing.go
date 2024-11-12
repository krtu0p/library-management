package dto

import (
	"base-gin/domain/dao"
	"time"
)

type BorrowingCreateReq struct {
	BorrowDate 	*time.Time `json:"borrow_date" binding:"omitempty"`
	ReturnDate 	*time.Time `json:"return_date" binding:"omitempty"`
	BookID    	uint   `json:"book_id" binding:"required"`
	PersonID 	uint   `json:"person_id" binding:"required"`
}

func (o *BorrowingCreateReq) ToEntity() dao.Borrowing {
	var item dao.Borrowing
	item.BorrowDate = o.BorrowDate 
	item.ReturnDate = o.ReturnDate 
	item.BookID = o.BookID
	item.PersonID = o.PersonID

	return item
}

type BorrowingResp struct {
	ID 				int 		`json:"id"`
	BorrowDate 		*time.Time 	`json:"borrow_date"`
	ReturnDate 		*time.Time 	`json:"return_date"`
	BorrowedBook   	string		`json:"borrowed_book"`
	BorrowerPerson 	string  	`json:"borrower_person"`
}

func (o *BorrowingResp) FromEntity(item *dao.Borrowing) {
	o.ID = int(item.ID)
	o.BorrowDate = item.BorrowDate
	o.ReturnDate = item.ReturnDate
	if item.BorrowedBook != nil {
        o.BorrowedBook = item.BorrowedBook.Title
    }
    if item.BorrowerPerson != nil {
        o.BorrowerPerson = item.BorrowerPerson.Fullname
    }
}

type BorrowingUpdateReq struct {
	ID 			uint 		`json:"-"`
	BorrowDate 	*time.Time 	`json:"borrow_date" binding:"omitempty"`
	ReturnDate 	*time.Time 	`json:"return_date" binding:"omitempty"`
	BookID    	uint   		`json:"book_id" binding:"required"`
	PersonID 	uint    	`json:"person_id" binding:"required"`
}
