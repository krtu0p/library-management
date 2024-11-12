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

func TestBook_Create_Success(t *testing.T) {
	p := dao.Publisher{
		Name: util.RandomStringAlpha(8),
		City: util.RandomStringAlpha(10),
	}
	_= publisherRepo.Create(&p)

	birthDate := time.Now().AddDate(-30, 0, 0)
	gender := domain.GenderMale
	a := dao.Author{
		Fullname: util.RandomStringAlpha(8),
		Gender: &gender,
		BirthDate: &birthDate,
	}
	_= authorRepo.Create(&a)
	
	params := dto.BookCreateReq{
		Title:       util.RandomStringAlpha(10),
		Subtitle:    util.RandomStringAlpha(15),
		AuthorID:    a.ID,
		PublisherID: p.ID,
	}
	
	w := doTest(
		"POST",
		server.RootBook, 
		params,
		createAuthAccessToken(dummyAdmin.Account.Username),
	)
	
	assert.Equal(t, 201, w.Code)

	fmt.Printf("%+v\n", params)
}

func TestBook_Update_Success(t *testing.T) {
	p := dao.Publisher{
		Name: util.RandomStringAlpha(8),
		City: util.RandomStringAlpha(10),
	}
	_= publisherRepo.Create(&p)

    birthDate := time.Now().AddDate(-30, 0, 0)
	gender := domain.GenderMale
	a := dao.Author{
		Fullname: util.RandomStringAlpha(8),
		Gender: &gender,
		BirthDate: &birthDate,
	}
	_= authorRepo.Create(&a)

    p2 := dao.Publisher{
        Name: util.RandomStringAlpha(7),
        City: util.RandomStringAlpha(8),
    }
    _ = publisherRepo.Create(&p2)

    a2 := dao.Author{
        Fullname:  util.RandomStringAlpha(9),
        Gender:    &gender,
        BirthDate: &birthDate,
    }
    _ = authorRepo.Create(&a2)

    b := dao.Book{
		Title:       util.RandomStringAlpha(10),
		Subtitle:    util.RandomStringAlpha(15),
		AuthorID:    a.ID,
		PublisherID: p.ID,
	}
	_ = bookRepo.Create(&b)

	fmt.Printf("First Borrowing: %+v\n", b)

	params := dto.BookUpdateReq{
		Title:       util.RandomStringAlpha(7),
		Subtitle:    util.RandomStringAlpha(12),
		AuthorID:    a2.ID,
		PublisherID: p2.ID,
	}

	w := doTest(
		"PUT",
		fmt.Sprintf("%s/%d", server.RootBook, b.ID),
		params,
		createAuthAccessToken(dummyAdmin.Account.Username),
	)
	assert.Equal(t, 200, w.Code)

	fmt.Printf("Updated Borrowing: %+v\n", params)
	item, _ := bookRepo.GetByID(b.ID)
	assert.Equal(t, params.Title, item.Title)
	assert.Equal(t, params.Subtitle, item.Subtitle)
	assert.Equal(t, params.AuthorID, item.AuthorID)
	assert.Equal(t, params.PublisherID, item.PublisherID)
}


func TestBook_GetList_Success(t *testing.T) {
    p := dao.Publisher{
        Name: util.RandomStringAlpha(8),
        City: util.RandomStringAlpha(10),
    }
    _ = publisherRepo.Create(&p)

    birthDate := time.Now().AddDate(-30, 0, 0)
    gender := domain.GenderMale
    a := dao.Author{
        Fullname:  util.RandomStringAlpha(8),
        Gender:    &gender,
        BirthDate: &birthDate,
    }
    _ = authorRepo.Create(&a)

    p2 := dao.Publisher{
        Name: util.RandomStringAlpha(8),
        City: util.RandomStringAlpha(10),
    }
    _ = publisherRepo.Create(&p2)

    a2 := dao.Author{
        Fullname:  util.RandomStringAlpha(8),
        Gender:    &gender,
        BirthDate: &birthDate,
    }
    _ = authorRepo.Create(&a2)

    b1 := dao.Book{
		Title:       util.RandomStringAlpha(10),
		Subtitle:    util.RandomStringAlpha(15),
		AuthorID:    a.ID,
		PublisherID: p.ID,
	}
	_ = bookRepo.Create(&b1)

    b2 := dao.Book{
		Title:       util.RandomStringAlpha(7),
		Subtitle:    util.RandomStringAlpha(14),
		AuthorID:    a2.ID,
		PublisherID: p2.ID,
	}
	_ = bookRepo.Create(&b2)

	w := doTest(
		"GET",
		server.RootBook,
		nil,
		"",
	)
	assert.Equal(t, 200, w.Code)

	body := w.Body.String()
	assert.Contains(t, body, b1.Title)
	assert.Contains(t, body, b2.Title)

	fmt.Printf("%+v\n", body)

	w = doTest(
		"GET",
		server.RootBook+"?q="+b1.Title,
		nil,
		"",
	)
	assert.Equal(t, 200, w.Code)

	body = w.Body.String() 
	assert.Contains(t, body, b1.Title)
	assert.NotContains(t, body, b2.Title)

	fmt.Printf("%+v\n", body)
}

func TestBook_GetByID_Success(t *testing.T) {
	p := dao.Publisher{
		Name: util.RandomStringAlpha(8),
		City: util.RandomStringAlpha(10),
	}
	_= publisherRepo.Create(&p)

	birthDate := time.Now().AddDate(-30, 0, 0)
	gender := domain.GenderMale
	a := dao.Author{
		Fullname: util.RandomStringAlpha(8),
		Gender: &gender,
		BirthDate: &birthDate,
	}
	_= authorRepo.Create(&a)

	b := dao.Book{
		Title:       util.RandomStringAlpha(10),
		Subtitle:    util.RandomStringAlpha(15),
		AuthorID:    a.ID,
		PublisherID: p.ID,
	}
	_ = bookRepo.Create(&b)

	o, err := bookRepo.GetByID(b.ID)
	assert.Nil(t, err)
	fmt.Printf("%+v\n", o)
}

func TestBook_Delete_Success(t *testing.T) {
	p := dao.Publisher{
		Name: util.RandomStringAlpha(8),
		City: util.RandomStringAlpha(10),
	}
	_= publisherRepo.Create(&p)

	birthDate := time.Now().AddDate(-30, 0, 0)
	gender := domain.GenderMale
	a := dao.Author{
		Fullname: util.RandomStringAlpha(8),
		Gender: &gender,
		BirthDate: &birthDate,
	}
	_= authorRepo.Create(&a)

	b := dao.Book{
		Title:       util.RandomStringAlpha(10),
		Subtitle:    util.RandomStringAlpha(15),
		AuthorID:    a.ID,
		PublisherID: p.ID,
	}
	_ = bookRepo.Create(&b)
	
	w := doTest(
		"DELETE",
		fmt.Sprintf("%s/%d", server.RootBook, b.ID),
		nil,
		createAuthAccessToken(dummyAdmin.Account.Username),
	)

	assert.Equal(t, 200, w.Code)

	item, _ := bookRepo.GetByID(b.ID)
	assert.Nil(t, item)
}