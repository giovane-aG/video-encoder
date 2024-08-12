package services_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/giovane-aG/video-encoder/encoder/application/repositories"
	"github.com/giovane-aG/video-encoder/encoder/application/services"
	"github.com/giovane-aG/video-encoder/encoder/domain"
	"github.com/giovane-aG/video-encoder/encoder/infrastructure/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatalf("error when loading env variables %s", err)
	}

	os.Setenv("CGO_ENABLED", "1")
}

func prepare() (repositories.VideoRepositoryInterface, *domain.Video) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewString()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repo := repositories.NewVideoRepository(db)

	return repo, video
}

func TestVideoServiceDownload(t *testing.T) {
	videoRepository, video := prepare()

	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = videoRepository

	err := videoService.Download("video-encoder-files")

	require.Nil(t, err)
}
