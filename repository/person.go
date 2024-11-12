package repository

import (
	"base-gin/domain/dao"
	"base-gin/domain/dto"
	"base-gin/exception"
	"base-gin/storage"
	"errors"


	"gorm.io/gorm"
)

type PersonRepository struct {
	db *gorm.DB
}

func NewPersonRepository(db *gorm.DB) *PersonRepository {
	return &PersonRepository{db: db}
}

func (r *PersonRepository) Create(newItem *dao.Person) error {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()

	tx := r.db.WithContext(ctx).Create(&newItem)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (r *PersonRepository) GetByAccountID(accountID uint) (dao.Person, error) {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()

	var item dao.Person
	tx := r.db.WithContext(ctx).Where(dao.Person{AccountID: &accountID}).
		First(&item)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return item, exception.ErrUserNotFound
		}

		return item, tx.Error
	}

	return item, nil
}

func (r *PersonRepository) GetByID(id uint) (*dao.Person, error) {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()

	var item dao.Person
	tx := r.db.WithContext(ctx).First(&item, id)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, exception.ErrUserNotFound
		}

		return nil, tx.Error
	}

	return &item, nil
}

func (r *PersonRepository) GetList(filter interface{}) ([]dao.Person, error) {
    // Cek jika filter nil
    if filter == nil {
        filter = &dao.Person{} // Default filter kosong
    }

    var persons []dao.Person
    result := r.db.Where(filter).Find(&persons)
    if result.Error != nil {
        return nil, result.Error
    }

    return persons, nil
}


func (r *PersonRepository) Update(params *dto.PersonUpdateReq) error {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()

	tx := r.db.WithContext(ctx).Model(&dao.Person{}).
		Where("id = ?", params.ID).
		Updates(map[string]interface{}{
			"fullname":   params.Fullname,
			"gender":     params.GetGender(),
			"birth_date": params.BirthDate,
		})

	return tx.Error
}

func (r *PersonRepository) Delete(id uint) error {
	ctx, cancelFunc := storage.NewDBContext()
	defer cancelFunc()

	tx := r.db.WithContext(ctx).Delete(&dao.Person{}, id)

	return tx.Error
}
