package services

import (
	"context"
	"io"
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

func (v *VideoService) Download(bucketName string) error {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	object := client.Bucket(bucketName).Object(v.Video.FilePath)

	reader, err := object.NewReader(ctx)
	if err != nil {
		return err
	}

	defer reader.Close()

	body, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	file, err := os.Create(
		os.Getenv("localStoragePath") + "/" + v.Video.ID + ".mp4")

	if err != nil {
		return err
	}

	_, err = file.Write(body)
	if err != nil {
		return err
	}

	defer file.Close()

	log.Printf("video %v has been stored", v.Video.ID)

	return nil
}

func (v *VideoService) Fragment() error {
	// create a new directory for the fragmented video
	err := os.Mkdir(os.Getenv("localStoragePath")+"/"+v.Video.ID, os.ModePerm)
	if err != nil {
		return err
	}

	sourceFile := os.Getenv("localStoragePath") + "/" + v.Video.ID + ".mp4"
	destinationFile := os.Getenv("localStoragePath") + "/" + v.Video.ID + ".frag"

	// fragments mp4 video
	command := exec.Command("mp4fragment", sourceFile, destinationFile)
	output, err := command.Output()
	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func (v *VideoService) Encode() error {
	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, os.Getenv("localStoragePath")+"/"+v.Video.ID+".frag")
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, os.Getenv("localStoragePath")+"/"+v.Video.ID)
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/opt/bento4/bin/")
	cmd := exec.Command("mp4dash", cmdArgs...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func (v *VideoService) Finish() error {
	err := os.Remove(os.Getenv("localStoragePath") + "/" + v.Video.ID + ".mp4")
	if err != nil {
		log.Println("error when trying to delete mp4" + v.Video.ID + ".mp4")
		return err
	}

	err = os.Remove(os.Getenv("localStoragePath") + "/" + v.Video.ID + ".frag")
	if err != nil {
		log.Println("error when trying to delete frag" + v.Video.ID + ".frag")
		return err
	}

	err = os.RemoveAll(os.Getenv("localStoragePath") + "/" + v.Video.ID)
	if err != nil {
		log.Println("error when trying to delete directory" + v.Video.ID)
		return err
	}

	log.Println("files have been removed")
	return nil
}

func (v *VideoService) InsertVideo() error {
	_, err := v.VideoRepository.Insert(v.Video)
	if err != nil {
		return err
	}
	return nil
}

func printOutput(out []byte) {
	log.Printf("Output ====> %s\n", string(out))
}
