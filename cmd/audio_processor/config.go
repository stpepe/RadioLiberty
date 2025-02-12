package audio_processor

import (
	"RadioLiberty/pkg/queue_service_client"
	"RadioLiberty/pkg/s3"
	"fmt"
	"log/slog"
	"os"
)

type GeneralConfig struct {
	Port             string
	DefaultAudioName string

	S3StorageConfig   *s3.Config
	QueueClientConfig *queue_service_client.Config
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

	nameDefaultAudioNameEnv := prefix + "DEFAULT_AUDIO_NAME"
	c.DefaultAudioName, exists = os.LookupEnv(nameDefaultAudioNameEnv)
	if !exists {
		c.DefaultAudioName = "default.mp3"
		slog.Warn("env variable %s is not set. Defaulting to %s", nameDefaultAudioNameEnv, c.DefaultAudioName)
	}

	c.S3StorageConfig = &s3.Config{}
	err := c.S3StorageConfig.EnvParse("S3_")
	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}

	c.QueueClientConfig = &queue_service_client.Config{}
	err = c.QueueClientConfig.EnvParse("QUEUE_CLIENT_")
	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}

	return nil
}
