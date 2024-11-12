package dto

import (
	"base-gin/domain"
	"base-gin/domain/dao"
	"time"
	"golang.org/x/crypto/bcrypt"
)

type AccountLoginReq struct {
	Username string `json:"uname" binding:"required,max=16"`
	Password string `json:"paswd" binding:"required,min=8,max=255"`
}

type AccountLoginResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AccountProfileResp struct {
	ID       int    `json:"id"`
	Fullname string `json:"fullname"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
}

func (o *AccountProfileResp) FromPerson(person *dao.Person) {
	var gender string
	if person.Gender == nil {
		gender = "-"
	} else if *person.Gender == domain.GenderFemale {
		gender = "wanita"
	} else {
		gender = "pria"
	}

	var age float64
	if person.BirthDate != nil {
		age = time.Since(*person.BirthDate).Hours() / (24 * 365)
	}

	o.Fullname = person.Fullname
	o.Gender = gender
	o.Age = int(age)
}


type AccountCreateReq struct {
	Username string `json:"uname" binding:"required,max=16"`
	Password string `json:"paswd" binding:"required,min=8,max=255"`
}

func (o *AccountCreateReq) ToEntity() dao.Account {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(o.Password), bcrypt.DefaultCost)


	return dao.Account{
		Username: o.Username,
		Password: string(hashedPassword),
	}
}

type AccountCreateResp struct{
	ID uint `json:"id"`
	Username string `json:"username"`
}

type AccountResp struct {
	ID uint `json:"id"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
}

type AccountUpdateReq struct {
	ID       uint   `json:"id" binding:"required"`
	Username string `json:"uname" binding:"required,max=16"`
	Password string `json:"paswd" binding:"required,min=8,max=255"`
	Email    string `json:"email" binding:"required,email"`
}

