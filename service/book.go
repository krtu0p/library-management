package service

import (
	"base-gin/domain/dto"
	"base-gin/exception"
	"base-gin/repository"
)

type BookService struct {
	repo *repository.BookRepository
}

func NewBookService(bookRepo *repository.BookRepository) *BookService {
	return &BookService{repo: bookRepo}
}

func (s *BookService) Create(params *dto.BookCreateReq) error {
	newItem := params.ToEntity()
	return s.repo.Create(&newItem)
}

func (s *BookService) GetByID(id uint) (dto.BookResp, error) {
	var resp dto.BookResp

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

func (s *BookService) GetList(params *dto.Filter) ([]dto.BookResp, error) {
	var resp []dto.BookResp

	items, err := s.repo.GetList(params)
	if err != nil {
		return nil, err
	}
	if len(items) < 1 {
		return nil, exception.ErrDataNotFound
	}

	for _, item := range items {
		var t dto.BookResp
		t.FromEntity(&item)

		resp = append(resp, t)
	}

	return resp, nil
}

func (s *BookService) Update(params *dto.BookUpdateReq) error {
	if params.ID <= 0 {
		return exception.ErrUserNotFound
	}

	return s.repo.Update(params)
}

func (s *BookService) Delete(id uint) error {
	if id <= 0 {
		return exception.ErrDataNotFound
	}

	return s.repo.Delete(id)
}