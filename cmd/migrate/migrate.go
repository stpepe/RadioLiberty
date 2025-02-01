package migrate

import (
	"RadioLiberty/pkg/local_storage"
	"RadioLiberty/pkg/s3"
	"context"
	"fmt"
	"mime/multipart"
	"os"
)

func Migrate() error {
	const errorMsg = "from migrate error: %w"

	config := &Config{}
	err := config.EnvParse("")
	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}

	s3Storage, err := s3.NewS3Storage(config.S3StorageConfig)
	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}

	audio, err := os.OpenFile(config.PathToDefaultAudioFile, os.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}
	defer audio.Close()
	info, err := audio.Stat()
	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}
	header := &multipart.FileHeader{
		Filename: info.Name(),
		Size:     info.Size(),
	}

	err = s3Storage.PutObject(context.Background(), audio, header)
	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}

	localStorage, err := local_storage.NewLocalStorage(config.LocalStorageConfig)
	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}

	err = localStorage.Migrate(config.PathToMigrations)
	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}

	return nil
}
