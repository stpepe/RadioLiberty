package s3

import (
	"RadioLiberty/pkg/models"
	"RadioLiberty/pkg/s3/minio_adapter"
	"context"
	"fmt"
	"maps"
	"mime/multipart"
	"slices"
)

var S3StorageRegister = map[string]func() (S3Storage, error){
	"minio": func() (S3Storage, error) { return minio_adapter.NewMinio() },
}

type S3Storage interface {
	PutObject(ctx context.Context, file multipart.File, header *multipart.FileHeader) error
	GetAudioFile(ctx context.Context, audioInfo *models.AudioInfo, removeFlag bool) (*models.AudioFile, error)
}

func NewS3Storage(config *Config) (S3Storage, error) {
	const errorMsg = "from NewS3Storage error: %w"
	const envVariableMustBeSetOneOf = "env variable %s must be set one of: %v"

	S3StorageConstructor, exists := S3StorageRegister[config.S3StorageType]
	if !exists {
		return nil, fmt.Errorf(
			errorMsg,
			fmt.Errorf(
				envVariableMustBeSetOneOf,
				config.S3StorageType,
				slices.Collect(maps.Keys(S3StorageRegister)),
			),
		)
	}

	S3Storage, err := S3StorageConstructor()
	if err != nil {
		return nil, fmt.Errorf(errorMsg, err)
	}

	return S3Storage, nil
}
