package repository

import (
	"base-gin/domain/dao"
	"base-gin/domain/dto"
	"base-gin/exception"
	"base-gin/storage"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) Create(newItem *dao.Book) error {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()

	tx := r.db.WithContext(ctx).Create(&newItem)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r *BookRepository) GetByID(id uint) (*dao.Book, error) {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()

	var item dao.Book
	tx := r.db.WithContext(ctx).
		Joins("BookPublisher").
		Joins("BookAuthor").
		First(&item, id)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, exception.ErrUserNotFound
		}
		return nil, tx.Error
	}
	return &item, nil
}

func (r *BookRepository) GetList(params *dto.Filter) ([]dao.Book, error) {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()

	var items []dao.Book
	tx := r.db.WithContext(ctx).
	Joins("BookPublisher").
	Joins("BookAuthor")

	if params.Keyword != "" {
		q := fmt.Sprintf("%%%s%%", params.Keyword)
		tx = tx.Where("title LIKE ?", q)
	}
	if params.Start >= 0 {
		tx = tx.Offset(params.Start)
	}
	if params.Limit > 0 {
		tx = tx.Limit(params.Limit)
	}

	tx = tx.Order("title ASC").Find(&items)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, tx.Error
	}

	return items, nil
}

func (r *BookRepository) Update(params *dto.BookUpdateReq) error {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()

    tx := r.db.WithContext(ctx).Model(&dao.Book{}).Where("id = ?", params.ID).Updates(map[string]interface{}{
        "title":       	params.Title,
        "subtitle":    	params.Subtitle,
        "author_id":    params.AuthorID,
        "publisher_id": params.PublisherID,
    })

	return tx.Error
}

func (r *BookRepository) Delete(id uint) error {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()

	tx := r.db.WithContext(ctx).Delete(&dao.Book{}, id)

	return tx.Error
}