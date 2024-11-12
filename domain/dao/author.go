package dao

import (
	"base-gin/domain"
	"time"
)

type Author struct {
	ID        uint 				 `gorm:"primarykey"`
	Fullname  string             `gorm:"size:56;not null;"`
	Gender    *domain.TypeGender `gorm:"type:enum('f','m');"`
	BirthDate *time.Time
}

func (Author) TableName() string {	
	return "authors"
}