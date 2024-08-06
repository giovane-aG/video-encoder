package repositories

import "github.com/giovane-aG/video-encoder/encoder/domain"

type JobRepositoryInterface interface {
	Insert(job *domain.Job) (*domain.Job, error)
	Find(id string) (*domain.Job, error)
	Update(job *domain.Job) (*domain.Job, error)
}
