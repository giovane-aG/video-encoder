package repositories_test

import (
	"testing"
	"time"

	"github.com/giovane-aG/video-encoder/encoder/application/repositories"
	"github.com/giovane-aG/video-encoder/encoder/domain"
	"github.com/giovane-aG/video-encoder/encoder/infrastructure/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestVideoRepositoryInsert(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewString()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repo := repositories.NewVideoRepository(db)
	_, err := repo.Insert(video)

	require.Nil(t, err)

	foundVideo, err := repo.Find(video.ID)

	require.Nil(t, err)
	require.NotEmpty(t, foundVideo.ID)
	require.Equal(t, video.ID, foundVideo.ID)
}
