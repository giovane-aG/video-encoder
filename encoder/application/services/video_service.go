package services

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"cloud.google.com/go/storage"

	"github.com/giovane-aG/video-encoder/encoder/application/repositories"
	"github.com/giovane-aG/video-encoder/encoder/domain"
)

type VideoService struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepositoryInterface
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (videoService *VideoService) Download(bucketName string) error {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	object := client.Bucket(bucketName).Object(videoService.Video.FilePath)

	reader, err := object.NewReader(ctx)
	if err != nil {
		return err
	}

	defer reader.Close()

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	file, err := os.Create(
		os.Getenv("localStoragePath") + "/" + videoService.Video.ID + ".mp4")

	if err != nil {
		return err
	}

	_, err = file.Write(body)
	if err != nil {
		return err
	}

	defer file.Close()

	log.Printf("video %v has been stored", videoService.Video.ID)

	return nil
}

func (videoService *VideoService) Fragment() error {
	// create a new directory for the fragmented video
	err := os.Mkdir(os.Getenv("localStoragePath")+"/"+videoService.Video.ID, os.ModePerm)
	if err != nil {
		return err
	}

	sourceFile := os.Getenv("localStoragePath" + "/" + videoService.Video.ID + ".mp4")
	destinationFile := os.Getenv("localStoragePath" + "/" + videoService.Video.ID + ".frag")

	// fragments mp4 video
	command := exec.Command("mp4fragment", sourceFile, destinationFile)
	output, err := command.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func printOutput(out []byte) {
	log.Printf("Output ====> %s\n", string(out))
}
