package migrate

import (
	"RadioLiberty/pkg/local_storage"
	"RadioLiberty/pkg/s3"
	"fmt"
	"os"
)

type Config struct {
	PathToMigrations       string
	PathToDefaultAudioFile string
	LocalStorageConfig     *local_storage.Config
	S3StorageConfig        *s3.Config
}

func (c *Config) EnvParse(prefix string) error {
	const errorMsg = "from general config error: %w"

	var exists bool

	namePathToMigrationsEnv := prefix + "PATH_TO_MIGRATIONS"
	c.PathToMigrations, exists = os.LookupEnv(namePathToMigrationsEnv)
	if !exists {
		c.PathToMigrations = "./migrations"
	}

	namePathToDefaultAudioFileEnv := prefix + "PATH_TO_DEFAULT_AUDIO_FILE"
	c.PathToDefaultAudioFile, exists = os.LookupEnv(namePathToDefaultAudioFileEnv)
	if !exists {
		c.PathToDefaultAudioFile = "./default_audio.mp3"
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
