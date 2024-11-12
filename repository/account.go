	package repository

	import (
		"base-gin/domain/dao"
		"base-gin/exception"
		"base-gin/storage"
		"errors"




		"gorm.io/gorm"
	)

	type AccountRepository struct {
		db *gorm.DB
	}

	func NewAccountRepository(db *gorm.DB) *AccountRepository {
		return &AccountRepository{db: db}
	}

	func (r *AccountRepository) Create(newItem *dao.Account) error {
		ctx, cancelFunc := storage.NewDBContext()
		defer cancelFunc()

		tx := r.db.WithContext(ctx).Create(&newItem)
		if tx.Error != nil {
			return tx.Error
		}

		return nil
	}

	func (r *AccountRepository) GetByUsername(uname string) (dao.Account, error) {
		ctx, cancelFunc := storage.NewDBContext()
		defer cancelFunc()

		var item dao.Account
		tx := r.db.WithContext(ctx).Where(dao.Account{Username: uname}).
			First(&item)
		if tx.Error != nil {
			if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
				return item, exception.ErrUserNotFound
			}

			return item, tx.Error
		}

		return item, nil
	}

	func (r *AccountRepository) Delete(id uint) error {
		// Check if there are any records in the persons table that reference the account
		var persons []dao.Person
		err := r.db.Where("account_id = ?", id).Find(&persons).Error
		if err != nil {
			return err
		}
		if len(persons) > 0 {
			// Delete the records in the persons table
			err := r.db.Where("account_id = ?", id).Delete(&dao.Person{}).Error
			if err != nil {
				return err
			}
		}
		// Delete the account
		err = r.db.Where("id = ?", id).Delete(&dao.Account{}).Error
		return err
	}

	func (r *AccountRepository) GetByID(id uint) (*dao.Account, error) {
		ctx, cancelFunc := storage.NewDBContext()
		defer cancelFunc()

		var account dao.Account
		tx := r.db.WithContext(ctx).First(&account, id)

		if tx.Error != nil && tx.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return &account, tx.Error
	}

	func (r *AccountRepository) GetByToken(token string) (*dao.Account, error) {
		ctx, cancelFunc := storage.NewDBContext()
		defer cancelFunc()

		var item dao.Account
		tx := r.db.WithContext(ctx).Where(dao.Account{Token: token}).
			First(&item)
		if tx.Error != nil {
			if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
				return nil, exception.ErrUserNotFound
			}

			return nil, tx.Error
		}

		return &item, nil
	}

	func (r *AccountRepository) GetList() ([]dao.Account, error) {
		ctx, cancelFunc := storage.NewDBContext()
		defer cancelFunc()

		var accounts []dao.Account
		tx := r.db.WithContext(ctx).Find(&accounts)

		return accounts, tx.Error
	}

	func (r *AccountRepository) Update(newItem *dao.Account) error {
		ctx, cancelFunc := storage.NewDBContext()
		defer cancelFunc()

		tx := r.db.WithContext(ctx).Model(&dao.Account{}).
			Where("id = ?", newItem.ID).
			Updates(map[string]interface{}{
				"username": newItem.Username,
				"password": newItem.Password,
			})

		return tx.Error
	}

		