package service

import (
	"base-gin/domain/dto"
	"base-gin/exception"
	"base-gin/repository"
)

type AuthorService struct {
	repo *repository.AuthorRepository
}

func NewAuthorService(authorRepo *repository.AuthorRepository) *AuthorService {
	return &AuthorService{repo: authorRepo}
}

func (s *AuthorService) Create(params *dto.AuthorCreateReq) error {
	newItem := params.ToEntity()
	return s.repo.Create(&newItem)
}

func (s *AuthorService) GetByID(id uint) (dto.AuthorResp, error) {
	var resp dto.AuthorResp

	item, err := s.repo.GetByID(id)
	if err != nil {
		return resp, err
	}
	if item == nil {
		return resp, exception.ErrDataNotFound
	}

	resp.FromEntity(item)

	return resp, nil
}

func (s *AuthorService) GetList(params *dto.Filter) ([]dto.AuthorResp, error) {
	var resp []dto.AuthorResp

	items, err := s.repo.GetList(params)
	if err != nil {
		return nil, err
	}
	if len(items) < 1 {
		return nil, exception.ErrDataNotFound
	}

	for _, item := range items {
		var t dto.AuthorResp
		t.FromEntity(&item)

		resp = append(resp, t)
	}

	return resp, nil
}

func (s *AuthorService) Update(params *dto.AuthorUpdateReq) error {
	if params.ID <= 0 {
		return exception.ErrDataNotFound
	}

	return s.repo.Update(params)
}

func (s *AuthorService) Delete(id uint) error {
	if id <= 0 {
		return exception.ErrDataNotFound
	}

	return s.repo.Delete(id)
}