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

func TestAuthor_Create_Success(t *testing.T) {
	birthDate := time.Now().AddDate(-30, 0, 0)
	gender := "f"
	params := dto.AuthorCreateReq{
		Fullname:  util.RandomStringAlpha(10),
		Gender:    gender,
		BirthDate: &birthDate,
	}

	w := doTest(
		"POST",
		server.RootAuthor,
		params,
		createAuthAccessToken(dummyAdmin.Account.Username),
	)
	
	assert.Equal(t, 201, w.Code)
	fmt.Printf("%+v\n", params)
}

func TestAuthor_Update_Success(t *testing.T) {
	birthDate := time.Now().AddDate(-30, 0, 0)
	gender := domain.GenderMale
	o := dao.Author{
		Fullname:  util.RandomStringAlpha(10),
		Gender:    &gender,
		BirthDate: &birthDate,
	}
	_ = authorRepo.Create(&o)

	newGender := domain.GenderFemale
	params := dto.AuthorUpdateReq{
		ID:       o.ID,
		Fullname: util.RandomStringAlpha(7),
		Gender:    &newGender,
		BirthDate: &birthDate,
	}

	w := doTest(
		"PUT",
		fmt.Sprintf("%s/%d", server.RootAuthor, o.ID),
		params,
		createAuthAccessToken(dummyAdmin.Account.Username),
	)
	assert.Equal(t, 200, w.Code)

	item, _ := authorRepo.GetByID(o.ID)
	assert.Equal(t, params.Fullname, item.Fullname)
	assert.Equal(t, params.Gender, item.Gender)
	assert.WithinDuration(t, *params.BirthDate, *item.BirthDate, time.Second, "Birth dates do not match")
}

func TestAuthor_Delete_Success(t *testing.T) {
	birthDate := time.Now().AddDate(-30, 0, 0)
	gender := domain.GenderMale
	o := dao.Author{
		Fullname:  util.RandomStringAlpha(10),
		Gender:    &gender,
		BirthDate: &birthDate,
	}
	_ = authorRepo.Create(&o)

	w := doTest(
		"DELETE",
		fmt.Sprintf("%s/%d", server.RootAuthor, o.ID),
		nil,
		createAuthAccessToken(dummyAdmin.Account.Username),
	)
	assert.Equal(t, 200, w.Code)

	item, _ := authorRepo.GetByID(o.ID)
	assert.Nil(t, item)
}

func TestAuthor_GetList_Success(t *testing.T) {
	birthDate := time.Now().AddDate(-30, 0, 0)
	gender := domain.GenderMale
	o1 := dao.Author{
		Fullname:  util.RandomStringAlpha(10),
		Gender:    &gender,
		BirthDate: &birthDate,
	}
	_ = authorRepo.Create(&o1)

	o2 := dao.Author{
		Fullname:  util.RandomStringAlpha(6),
		Gender:    &gender,
		BirthDate: &birthDate,
	}
	_ = authorRepo.Create(&o2)

	w := doTest(
		"GET",
		server.RootAuthor,
		nil,
		"",
	)
	assert.Equal(t, 200, w.Code)

	body := w.Body.String()
	assert.Contains(t, body, o1.Fullname)
	assert.Contains(t, body, o2.Fullname)

	w = doTest(
		"GET",
		server.RootAuthor+"?q="+o1.Fullname,
		nil,
		"",
	)
	assert.Equal(t, 200, w.Code)

	body = w.Body.String()
	assert.Contains(t, body, o1.Fullname)
	assert.NotContains(t, body, o2.Fullname)
}

func TestAuthor_GetDetail_Success(t *testing.T) {
	birthDate := time.Now().AddDate(-30, 0, 0)
	gender := domain.GenderMale
	o := dao.Author{
		Fullname:  util.RandomStringAlpha(10),
		Gender:    &gender,
		BirthDate: &birthDate,
	}
	_ = authorRepo.Create(&o)

	w := doTest(
		"GET",
		fmt.Sprintf("%s/%d", server.RootAuthor, o.ID),
		nil,
		"",
	)
	assert.Equal(t, 200, w.Code)

	body := w.Body.String()
	assert.Contains(t, body, o.Fullname)
}
