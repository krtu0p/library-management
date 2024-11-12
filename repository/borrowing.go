package repository

import (
	"base-gin/domain/dao"
	"base-gin/domain/dto"
	"base-gin/exception"
	"base-gin/storage"
	"errors"

	"gorm.io/gorm"
)

type BorrowingRepository struct {
	db *gorm.DB
}

func NewBorrowingRepository(db *gorm.DB) *BorrowingRepository {
	return &BorrowingRepository{db: db}
}

func (r *BorrowingRepository) Create(newItem *dao.Borrowing) error {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()

	tx := r.db.WithContext(ctx).Create(&newItem)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r *BorrowingRepository) GetByID(id uint) (*dao.Borrowing, error) {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()

	var item dao.Borrowing
	tx := r.db.WithContext(ctx).
		Joins("BorrowedBook").
		Joins("BorrowerPerson").
		First(&item, id)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, exception.ErrUserNotFound
		}
		return nil, tx.Error
	}
	return &item, nil
}

func (r *BorrowingRepository) GetList(params *dto.Filter) ([]dao.Borrowing, error) {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()

	var items []dao.Borrowing
	tx := r.db.WithContext(ctx).
    Joins("BorrowedBook").
    Joins("BorrowerPerson")

	if params.Start >= 0 {
		tx = tx.Offset(params.Start)
	}
	if params.Limit > 0 {
		tx = tx.Limit(params.Limit)
	}

	tx = tx.Find(&items)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, tx.Error
	}

	return items, nil
}

func (r *BorrowingRepository) Update(params *dto.BorrowingUpdateReq) error {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()

    tx := r.db.WithContext(ctx).Joins("BorrowedBook").Joins("BorrowerPerson").Model(&dao.Borrowing{}).Where("id = ?", params.ID).Updates(map[string]interface{}{
        "borrow_date":	params.BorrowDate,
        "return_date": 	params.ReturnDate,
        "book_id":    params.BookID,
        "person_id": params.PersonID,
    })

	return tx.Error
}

func (r *BorrowingRepository) Delete(id uint) error {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()

	tx := r.db.WithContext(ctx).Delete(&dao.Borrowing{}, id)

	return tx.Error
}