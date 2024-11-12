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

type AuthorRepository struct {
	db *gorm.DB
}

func NewAuthorRepository(db *gorm.DB) *AuthorRepository {
	return &AuthorRepository{db: db}
}

func (r *AuthorRepository) Create(newItem *dao.Author) error {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()

	tx := r.db.WithContext(ctx).Create(&newItem)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (r *AuthorRepository) GetByID(id uint) (*dao.Author, error) {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()

	var item dao.Author
	tx := r.db.WithContext(ctx).First(&item, id)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, exception.ErrUserNotFound
		}

		return nil, tx.Error
	}

	return &item, nil
}

func (r *AuthorRepository) GetList(params *dto.Filter) ([]dao.Author, error) {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()

	var items []dao.Author
	tx := r.db.WithContext(ctx)

	if params.Keyword != "" {
		q := fmt.Sprintf("%%%s%%", params.Keyword)
		tx = tx.Where("fullname LIKE ?", q)
	}
	if params.Start >= 0 {
		tx = tx.Offset(params.Start)
	}
	if params.Limit > 0 {
		tx = tx.Limit(params.Limit)
	}

	tx = tx.Order("fullname ASC").Find(&items)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, tx.Error
	}

	return items, nil
}

func (r *AuthorRepository) Update(params *dto.AuthorUpdateReq) error {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()

	tx := r.db.WithContext(ctx).Model(&dao.Author{}).
		Where("id = ?", params.ID).
		Updates(map[string]interface{}{
			"fullname": params.Fullname,
			"gender": params.Gender,
			"birth_date": params.BirthDate,
		})

	return tx.Error
}

func (r *AuthorRepository) Delete(id uint) error {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()

	tx := r.db.WithContext(ctx).Delete(&dao.Author{}, id)

	return tx.Error
}
