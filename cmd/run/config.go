package run

import (
	"RadioLiberty/pkg/local_storage"
	"RadioLiberty/pkg/s3"
	"fmt"
	"log/slog"
	"os"
)

type GeneralConfig struct {
	Port string

	LocalStorageConfig *local_storage.Config
	S3StorageConfig    *s3.Config
}

func (c *GeneralConfig) EnvParse(prefix string) error {
	const errorMsg = "from general config error: %w"

	var exists bool

	namePortEnv := prefix + "PORT"
	c.Port, exists = os.LookupEnv(namePortEnv)
	if !exists {
		c.Port = "8080"
		slog.Warn("env variable %s is not set. Defaulting to %s", namePortEnv, c.Port)
	}

	c.LocalStorageConfig = &local_storage.Config{}
	err := c.LocalStorageConfig.EnvParse("LOCAL_STORAGE_")
	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}

	c.S3StorageConfig = &s3.Config{}
	err = c.S3StorageConfig.EnvParse("S3_")
	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}

	return nil
}
