package repository

import (
	"go-url-shortener/internal/model"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	DB *gorm.DB
}

func NewPostgresRepo(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{DB: db}
}

func (r *PostgresRepository) CreateURL(url *model.URL) error {
	result := r.DB.Create(url)
	
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *PostgresRepository) FindByShortCode(shortCode string) (*model.URL, error) {
	var url model.URL
	// Use GORM's First method to find the record based on the ShortCode field.
	result := r.DB.Where("short_code = ?", shortCode).First(&url)
	
	if result.Error != nil {
		// If no record is found, GORM returns gorm.ErrRecordNotFound
		return nil, result.Error
	}
	return &url, nil
}