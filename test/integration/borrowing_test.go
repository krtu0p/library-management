package integration_test

import (
	"base-gin/domain"
	"base-gin/domain/dao"
	"base-gin/domain/dto"
	"base-gin/server"
	"base-gin/util"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func CreateAuthor()*dao.Author {
	birthDate := time.Now().AddDate(-30, 0, 0)
	gender := domain.GenderMale
	a := dao.Author{
		Fullname:  util.RandomStringAlpha(8),
		Gender:    &gender,
		BirthDate: &birthDate,
	}

	db.Create(&a)

	return &a
}

func CreatePublisher() *dao.Publisher {
	p := dao.Publisher{
		Name: util.RandomStringAlpha(8),
		City: util.RandomStringAlpha(10),
	}

	db.Create(&p)

	return &p
}

func CreateBook() *dao.Book {
	a := CreateAuthor()
	p := CreatePublisher()

	b := dao.Book{
		Title:       util.RandomStringAlpha(10),
		Subtitle:    util.RandomStringAlpha(15),
		AuthorID:    a.ID,                      
		PublisherID: p.ID,                     
	}

	db.Create(&b)

	return &b
}

func CreatePerson() *dao.Person {
	birthDate := time.Now().AddDate(-30, 0, 0)
	gender := domain.GenderFemale
	p := dao.Person{
		Fullname:  util.RandomStringAlpha(8) + " " + util.RandomStringAlpha(8),
		Gender:    &gender,
		BirthDate: &birthDate,
	}
	db.Create(&p)
	return &p
}

func TestBorrowing_Create_Success(t *testing.T) {
	b := CreateBook()
	p := CreatePerson()

	borrowDate := time.Now()
	returnDate := borrowDate.AddDate(0, 0, 7) 
	params := dto.BorrowingCreateReq{
		BorrowDate: &borrowDate,
		ReturnDate: &returnDate,
		BookID:     b.ID,
		PersonID:   p.ID,
	}

	w := doTest(
		"POST",
		server.RootBorrowing,
		params,
		createAuthAccessToken(dummyAdmin.Account.Username),
	)
	
	assert.Equal(t, 201, w.Code)
	
	assert.NotNil(t, params.BorrowDate)
	assert.NotNil(t, params.ReturnDate)
	assert.Equal(t, b.ID, params.BookID)
	assert.Equal(t, p.ID, params.PersonID)
	
	fmt.Printf("%+v\n", params)
}

func TestBorrowing_Update_Success(t *testing.T) {
	b := CreateBook()
	b2 := CreateBook()
	p := CreatePerson()

	borrowDate := time.Now()
	returnDate := borrowDate.AddDate(0, 0, 7) 
	params := dao.Borrowing{
		BorrowDate: &borrowDate,
		ReturnDate: &returnDate,
		BookID:     b.ID,
		PersonID:   p.ID,
	}
	_ = borrowingRepo.Create(&params)

	fmt.Printf("First Borrowing: %+v\n", params)
	paramsUpdate := dto.BorrowingUpdateReq{
		BorrowDate: &borrowDate,
		ReturnDate: &returnDate,
		BookID:     b2.ID,
		PersonID:   p.ID,
	}

	w := doTest(
		"PUT",
		fmt.Sprintf("%s/%d", server.RootBorrowing, params.ID),
		paramsUpdate,
		createAuthAccessToken(dummyAdmin.Account.Username),
	)

	assert.Equal(t, 200, w.Code)

	// Fokus perhatikan pada BookID
	fmt.Printf("Updated Borrowing: %+v\n", paramsUpdate)

	item, _ := borrowingRepo.GetByID(params.ID)
	assert.WithinDuration(t, *paramsUpdate.BorrowDate, *item.BorrowDate, time.Second)
	assert.WithinDuration(t, *paramsUpdate.ReturnDate, *item.ReturnDate, time.Second)
	assert.Equal(t, paramsUpdate.BookID, item.BookID)
	assert.Equal(t, paramsUpdate.PersonID, item.PersonID)
}

func TestBorrowing_Getlist_Success(t *testing.T) {
	b1 := CreateBook()
	b2 := CreateBook()
	p1 := CreatePerson()
	p2 := CreatePerson()

	borrowDate := time.Now()
	returnDate := borrowDate.AddDate(0, 0, 7) 
	params1 := dao.Borrowing{
		BorrowDate: &borrowDate,
		ReturnDate: &returnDate,
		BookID:     b1.ID,
		PersonID:   p1.ID,
	}
	_ = borrowingRepo.Create(&params1)

	params2 := dao.Borrowing{
		BorrowDate: &borrowDate,
		ReturnDate: &returnDate,
		BookID:     b2.ID,
		PersonID:   p2.ID,
	}
	_ = borrowingRepo.Create(&params2)

	w := doTest(
		"GET",
		server.RootBorrowing,
		nil,
		"",
	)

	assert.Equal(t, 200, w.Code)
	
	body := w.Body.String()
	fmt.Printf("%+v\n", body)
}

func TestBorrowing_GetByID_Success(t *testing.T) {
	b := CreateBook()
	p := CreatePerson()
	
	borrowDate := time.Now()
	returnDate := borrowDate.AddDate(0, 0, 7) 
	params := dao.Borrowing{
		BorrowDate: &borrowDate,
		ReturnDate: &returnDate,
		BookID:     b.ID,
		PersonID:   p.ID,
	}
	_ = borrowingRepo.Create(&params)

	o, err := bookRepo.GetByID(params.ID)
	assert.Nil(t, err)
	fmt.Printf("%+v\n", o)
}

func TestBorrowing_Delete_Success(t *testing.T) {
	b := CreateBook()
	p := CreatePerson()

	borrowDate := time.Now()
	returnDate := borrowDate.AddDate(0, 0, 7) 
	params := dao.Borrowing{
		BorrowDate: &borrowDate,
		ReturnDate: &returnDate,
		BookID:     b.ID,
		PersonID:   p.ID,
	}
	_ = borrowingRepo.Create(&params)

	w := doTest(
		"DELETE",
		fmt.Sprintf("%s/%d", server.RootBorrowing, params.ID),
		nil,
		createAuthAccessToken(dummyAdmin.Account.Username),
	)
	assert.Equal(t, 200, w.Code)

	item, _ := borrowingRepo.GetByID(params.ID)
	assert.Nil(t, item)
}