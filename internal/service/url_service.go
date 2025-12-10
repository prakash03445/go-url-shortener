package service

import (
	"go-url-shortener/internal/model"
	"go-url-shortener/pkg/utils"
)

type PostgresRepo interface {
	CreateURL(url *model.URL) error
}

type URLService struct {
	PostgresRepo PostgresRepo
}

func NewURLService(pgRepo PostgresRepo) *URLService {
	return &URLService{PostgresRepo: pgRepo}
}

func (s *URLService) Shorten(longURL string) (*model.URL, error) {
	shortCode := utils.GenerateShortCode(7) 

	url := &model.URL{
		ShortCode: shortCode,
		LongURL:   longURL,
	}

	if err := s.PostgresRepo.CreateURL(url); err != nil {
		return nil, err
	}

	return url, nil
}