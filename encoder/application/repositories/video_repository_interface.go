package repositories

import "github.com/giovane-aG/video-encoder/encoder/domain"

type VideoRepositoryInterface interface {
	Insert(video *domain.Video) (*domain.Video, error)
	Find(id string) (*domain.Video, error)
}
