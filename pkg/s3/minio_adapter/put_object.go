package minio_adapter

import (
	"RadioLiberty/pkg/models"
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
)

func (ma *MinioAdapter) PutObject(ctx context.Context, file multipart.File, header *multipart.FileHeader) error {
	const errorMsg = "from minio PutObject error: %w"

	exists, err := ma.client.BucketExists(ctx, ma.minioBucket)
	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}
	if !exists {
		err = ma.client.MakeBucket(ctx, ma.minioBucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf(errorMsg, err)
		}
	}

	_, err = ma.client.PutObject(ctx, ma.minioBucket, header.Filename, file, header.Size, minio.PutObjectOptions{
		ContentType: header.Header.Get("Content-Type"),
	})
	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}
	return nil
}

func (ma *MinioAdapter) GetObject(ctx context.Context, objectName string) (*minio.Object, error) {
	const errorMsg = "from minio GetObject error: %w"

	audioFile, err := ma.client.GetObject(ctx, ma.minioBucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf(errorMsg, err)
	}
	return audioFile, nil
}

func (ma *MinioAdapter) GetAudioFile(ctx context.Context, audioInfo *models.AudioInfo, removeFlag bool) (*models.AudioFile, error) {
	const errorMsg = "from minio GetAudioFile error: %w"

	audioFile := &models.AudioFile{
		Info: *audioInfo,
	}
	file, err := ma.client.GetObject(ctx, ma.minioBucket, audioInfo.FileName, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf(errorMsg, err)
	}
	defer file.Close()
	audioFile.File, err = io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf(errorMsg, err)
	}

	if removeFlag {
		err = ma.client.RemoveObject(ctx, ma.minioBucket, audioInfo.FileName, minio.RemoveObjectOptions{})
		if err != nil {
			return nil, fmt.Errorf(errorMsg, err)
		}
	}
	return audioFile, nil
}
