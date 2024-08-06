package repositories

import (
	"github.com/giovane-aG/video-encoder/encoder/domain"
	"github.com/jinzhu/gorm"
)

type JobRepository struct {
	DB *gorm.DB
}

func (repo JobRepository) Insert(job *domain.Job) (*domain.Job, error) {

	err := repo.DB.Create(job).Error

	if err != nil {
		return nil, err
	}

	return job, nil
}
