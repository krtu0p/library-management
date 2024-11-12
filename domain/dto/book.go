package dto

import "base-gin/domain/dao"

type BookCreateReq struct {
	Title       string `json:"title" binding:"required,max=56"`
	Subtitle    string `json:"subtitle" binding:"required,max=64"`
	AuthorID    uint   `json:"author_id" binding:"required"`
	PublisherID uint   `json:"publisher_id" binding:"required"`
}

func (o *BookCreateReq) ToEntity() dao.Book {
	var item dao.Book
	item.Title = o.Title 
	item.Subtitle = o.Subtitle 
	item.AuthorID = o.AuthorID
	item.PublisherID = o.PublisherID

	return item
}

type BookResp struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
	Author    string `json:"author"`
	Publisher string `json:"publisher"`
}

func (o *BookResp) FromEntity(item *dao.Book) {
	o.ID = int(item.ID)
	o.Title = item.Title
	o.Subtitle = item.Subtitle
	if item.BookAuthor != nil {
        o.Author = item.BookAuthor.Fullname
    }
    if item.BookPublisher != nil {
        o.Publisher = item.BookPublisher.Name
    }
}

type BookUpdateReq struct {
	ID          uint   `json:"-"`
	Title       string `json:"title" binding:"required,max=56"`
	Subtitle    string `json:"subtitle" binding:"required,max=64"`
	AuthorID    uint   `json:"author_id" binding:"required"`
	PublisherID uint   `json:"publisher_id" binding:"required"`
}
