package service

import (
    "errors"
	"log"

	"go-url-shortener/internal/model"
	"go-url-shortener/pkg/utils"
	
	"gorm.io/gorm"
)

type PostgresRepo interface {
	CreateURL(url *model.URL) error
	FindByShortCode(shortCode string) (*model.URL, error)
}

type RedisRepo interface {
	GetLongURL(shortCode string) (string, error)
	SetLongURL(shortCode string, longURL string) error
}

var ErrURLNotFound = errors.New("short URL not found")

type URLService struct {
	PostgresRepo PostgresRepo
	RedisRepo    RedisRepo
}

func NewURLService(pgRepo PostgresRepo, redisRepo RedisRepo) *URLService {
	return &URLService{PostgresRepo: pgRepo, RedisRepo: redisRepo}
}

func (s *URLService) ResolveURL(shortCode string) (string, error) {
	longURL, err := s.RedisRepo.GetLongURL(shortCode)
	if err != nil {
		log.Printf("Warning: Redis read error, treating as cache miss: %v", err)
	}
	if longURL != "" {
		return longURL, nil
	}

	urlModel, err := s.PostgresRepo.FindByShortCode(shortCode) 
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { 
			return "", ErrURLNotFound 
		}
		return "", err
	}

	if err := s.RedisRepo.SetLongURL(shortCode, urlModel.LongURL); err != nil {
		log.Printf("Warning: Failed to populate cache for code %s: %v", shortCode, err)
	}

	return urlModel.LongURL, nil
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