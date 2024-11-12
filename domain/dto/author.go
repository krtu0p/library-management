package dto

import (
	"base-gin/domain"
	"base-gin/domain/dao"
	"time"
)

type AuthorCreateReq struct {
	Fullname  string 	 `json:"fullname" binding:"required,max=56"`
	Gender    string     `json:"gender" binding:"required,oneof=m f"`
	BirthDate *time.Time `json:"birth_date" binding:"omitempty"`
}


func (o *AuthorCreateReq) ToEntity() dao.Author {
	var item dao.Author
	item.Fullname = o.Fullname

	var gender domain.TypeGender
	if o.Gender == "m" {
		gender = domain.GenderMale
	} else if o.Gender == "f" {
		gender = domain.GenderFemale
	}
	item.Gender = &gender

	item.BirthDate = o.BirthDate

	return item
}

type AuthorResp struct {
	ID 		  int 					`json:"id"`
	Fullname  string    			`json:"fullname"`
	Gender    *domain.TypeGender    `json:"gender"`
	BirthDate *time.Time 			`json:"birth_date"`
}

func (o *AuthorResp) FromEntity(item *dao.Author) {
	o.ID = int(item.ID)
	o.Fullname = item.Fullname
	o.Gender = item.Gender
	o.BirthDate = item.BirthDate
}

type AuthorUpdateReq struct {
	ID   	  uint   				`json:"-"`
	Fullname  string    			`json:"fullname" binding:"required,max=56"`
	Gender    *domain.TypeGender    `json:"gender" binding:"omitempty,oneof=f m"`
	BirthDate *time.Time 			`json:"birth_date" binding:"omitempty"`
}
