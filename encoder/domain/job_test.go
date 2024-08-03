package domain_test

import (
	"testing"
	"time"

	"github.com/giovane-aG/video-encoder/encoder/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestValidateJob(t *testing.T) {
	video := domain.NewVideo()
	video.ResourceID = "1234"
	video.FilePath = "filePath"
	video.ID = uuid.NewString()
	video.CreatedAt = time.Now()

	job, err := domain.NewJob("path", "Converted", video)

	require.NotNil(t, job)
	require.Nil(t, err)
}
