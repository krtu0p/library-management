package service

import (
	"base-gin/config"
	"base-gin/domain/dto"
	"base-gin/exception"
	"base-gin/repository"
	"base-gin/util"
	"base-gin/domain/dao"
	"errors"


	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

type AccountService struct {
	cfg  *config.Config
	repo *repository.AccountRepository
}

func NewAccountService(
	cfg *config.Config,
	accountRepo *repository.AccountRepository,
) *AccountService {
	return &AccountService{cfg: cfg, repo: accountRepo}
}



func (s *AccountService) Login(p dto.AccountLoginReq) (dto.AccountLoginResp, error) {
	var resp dto.AccountLoginResp

	item, err := s.repo.GetByUsername(p.Username)
	if err != nil {
		return resp, err
	}

	if paswdOk := item.VerifyPassword(p.Password); !paswdOk {
		return resp, exception.ErrUserLoginFailed
	}

	aToken, err := util.CreateAuthAccessToken(*s.cfg, item.Username)
	if err != nil {
		return resp, err
	}

	rToken, err := util.CreateAuthRefreshToken(*s.cfg, item.Username)
	if err != nil {
		return resp, err
	}

	resp.AccessToken = aToken
	resp.RefreshToken = rToken

	return resp, nil
}

func (s *AccountService) Create(params dto.AccountCreateReq) (*dao.Account, error) {
	newItem := params.ToEntity()
	err := s.repo.Create(&newItem)
	return &newItem, err
}


func (s *AccountService) Delete(id uint) error {
	if id <= 0 {
		return exception.ErrDataNotFound
	}

	return s.repo.Delete(id)
}

func (s *AccountService) GetByID(id uint) (*dao.Account, error) {
	return s.repo.GetByID(id)
}

func (s *AccountService) GetAccountByID(id uint) (*dao.Account, error) {
	// retrieve account by ID from database or repository
	account, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrDataNotFound
		}
		return nil, err
	}
	return account, nil
}

func (s *AccountService) GetAccountByToken(c *gin.Context) (*dao.Account, error) {
	// implement the logic to retrieve the account by token
	// for example:
	token := c.GetHeader("Authorization")
	account, err := s.repo.GetByToken(token)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *AccountService) GetList(params *dto.Filter) ([]dao.Account, error) {
	return s.repo.GetList()
}

func (s *AccountService) Update(params *dto.AccountUpdateReq) (dao.Account, error) {
	account := &dao.Account{
		ID:        params.ID,
		Username:  params.Username,
		Password:  params.Password,
	}
	return *account, s.repo.Update(account)
}