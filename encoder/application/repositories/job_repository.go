package repositories

import (
	"fmt"

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

func (repo JobRepository) Find(id string) (*domain.Job, error) {
	var job *domain.Job
	repo.DB.First(job, "id = ?", id)

	if job.ID == "" {
		return nil, fmt.Errorf("no job with this id was found")
	}

	return job, nil
}
