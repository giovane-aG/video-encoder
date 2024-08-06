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

func TestJobRepositoryDBInsert(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	jobRepository := repositories.NewJobRepository(db)
	videoRepository := repositories.NewVideoRepository(db)

	video := domain.NewVideo()
	video.ID = uuid.NewString()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	videoRepository.Insert(video)

	job, err := domain.NewJob(
		"output",
		"uploaded",
		video,
	)
	require.Nil(t, err)

	_, err = jobRepository.Insert(job)
	require.Nil(t, err)

	foundJob, err := jobRepository.Find(job.ID)
	require.Nil(t, err)
	require.Equal(t, foundJob.ID, job.ID)
	require.Equal(t, job.VideoID, video.ID)
}

func TestJobRepositoryDBUpdate(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	jobRepository := repositories.NewJobRepository(db)
	videoRepository := repositories.NewVideoRepository(db)

	video := domain.NewVideo()
	video.ID = uuid.NewString()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	videoRepository.Insert(video)

	job, err := domain.NewJob(
		"output",
		"uploaded",
		video,
	)
	require.Nil(t, err)

	_, err = jobRepository.Insert(job)
	require.Nil(t, err)

	job.OutputBucketPath = "new_output_path"
	_, err = jobRepository.Update(job)
	require.Nil(t, err)

	foundJob, err := jobRepository.Find(job.ID)
	require.Nil(t, err)
	require.Equal(t, foundJob.ID, job.ID)
	require.Equal(t, foundJob.OutputBucketPath, "new_output_path")
	require.Equal(t, job.VideoID, video.ID)
}
