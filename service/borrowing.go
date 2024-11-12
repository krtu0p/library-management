package service

import (
	"base-gin/domain/dto"
	"base-gin/exception"
	"base-gin/repository"
)

type BorrowingService struct {
	repo *repository.BorrowingRepository
}

func NewBorrowingService(borrowingRepo *repository.BorrowingRepository) *BorrowingService {
	return &BorrowingService{repo: borrowingRepo}
}

func (s *BorrowingService) Create(params *dto.BorrowingCreateReq) error {
	newItem := params.ToEntity()
	return s.repo.Create(&newItem)
}

func (s *BorrowingService) GetByID(id uint) (dto.BorrowingResp, error) {
	var resp dto.BorrowingResp

	item, err := s.repo.GetByID(id)
	if err != nil {
		return resp, err
	}
	if item == nil {
		return resp, exception.ErrUserNotFound
	}

	resp.FromEntity(item)

	return resp, nil
}

func (s *BorrowingService) GetList(params *dto.Filter) ([]dto.BorrowingResp, error) {
	var resp []dto.BorrowingResp

	items, err := s.repo.GetList(params)
	if err != nil {
		return nil, err
	}
	if len(items) < 1 {
		return nil, exception.ErrDataNotFound
	}

	for _, item := range items {
		var t dto.BorrowingResp
		t.FromEntity(&item)

		resp = append(resp, t)
	}

	return resp, nil
}

func (s *BorrowingService) Update(params *dto.BorrowingUpdateReq) error {
	if params.ID <= 0 {
		return exception.ErrUserNotFound
	}

	return s.repo.Update(params)
}

func (s *BorrowingService) Delete(id uint) error {
	if id <= 0 {
		return exception.ErrDataNotFound
	}

	return s.repo.Delete(id)
}