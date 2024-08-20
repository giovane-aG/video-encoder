package services

import (
	"context"
	"io"
	"os"
	"strings"

	"cloud.google.com/go/storage"
)

type VideoUploadManagerService struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Errors       []string
}

func NewVideoUploadManagerService() *VideoUploadManagerService {
	return &VideoUploadManagerService{}
}

func (vu *VideoUploadManagerService) UploadObject(objectPath string, client *storage.Client, ctx context.Context) error {
	path := strings.Split(objectPath, os.Getenv("localStoragePath"+"/"))

	file, err := os.Open(objectPath)
	if err != nil {
		return err
	}

	defer file.Close()

	wc := client.Bucket(vu.OutputBucket).Object(path[1]).NewWriter(ctx)
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	if _, err = io.Copy(wc, file); err != nil {
		return err
	}

	return nil
}

func getClientUpload() (*storage.Client, context.Context, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	return client, ctx, nil
}
