package services

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
	"io"
	"log"
	"mime/multipart"
	"telego/app/config"
)

type GoogleStorage struct {
	client *storage.Client
	bucket *storage.BucketHandle
}

func NewGoogleStorage(ctx context.Context) (*GoogleStorage, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	bucket := client.Bucket(config.Config.GCloud.GcloudStorageBucket)
	//defer client.Close()
	return &GoogleStorage{
		client,
		bucket,
	}, nil
}

func (gs GoogleStorage) UploadFile(ctx context.Context, f multipart.File, projectId string) (*uuid.UUID, error) {
	objectName := uuid.New()
	wc := gs.bucket.Object(projectId + "/" + objectName.String()).NewWriter(ctx)
	wc.ContentEncoding = "gzip"
	wc.CacheControl = "public, max-age=172800" // cache for 2 days

	if _, err := io.Copy(wc, f); err != nil {
		return nil, fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return nil, fmt.Errorf("Writer.Close: %v", err)
	}
	return &objectName, nil
}

func (gs GoogleStorage) CopyFiles(
	ctx context.Context,
	srcProjectId string,
	oldFilePaths []string,
	destProjectId string,
	newFilePaths []string,
) ([]*storage.ObjectAttrs, error) {
	var attrs []*storage.ObjectAttrs
	for index, oldFilePath := range oldFilePaths {
		src := gs.bucket.Object(srcProjectId + "/" + oldFilePath)
		dest := gs.bucket.Object(destProjectId + "/" + newFilePaths[index])

		attr, err := dest.CopierFrom(src).Run(ctx)
		if err != nil {
			log.Println(err)
		}
		attrs = append(attrs, attr)
	}
	return attrs, nil
}

func (gs GoogleStorage) DeleteFiles(ctx context.Context, filePaths []string) error {
	for _, filePath := range filePaths {
		if err := gs.bucket.Object(filePath).Delete(ctx); err != nil {
			log.Println(err)
		}

	}
	return nil
}

func (gs GoogleStorage) DeleteFolders(ctx context.Context, folderPath string) error {
	filesIter := gs.bucket.Objects(ctx, &storage.Query{Prefix: folderPath})
	for {
		objAttrs, err := filesIter.Next()
		if err != nil && err != iterator.Done {
			log.Println(err)
			return err
		}
		if err == iterator.Done {
			break
		}
		if err := gs.bucket.Object(objAttrs.Name).Delete(ctx); err != nil && !errors.Is(err, storage.ErrObjectNotExist) {
			log.Println(err)
		}
	}
	return nil
}
