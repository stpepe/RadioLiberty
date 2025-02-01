package minio_adapter

import (
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioAdapter struct {
	client      *minio.Client
	minioBucket string
}

func NewMinio() (*MinioAdapter, error) {
	const errorMsg = "from NewMinio error: %w"

	client, err := minio.New(MinioConfig.Address, &minio.Options{
		Creds:  credentials.NewStaticV4(MinioConfig.AccessKey, MinioConfig.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf(errorMsg, err)
	}

	adapter := &MinioAdapter{
		client:      client,
		minioBucket: MinioConfig.Bucket,
	}
	return adapter, nil
}
