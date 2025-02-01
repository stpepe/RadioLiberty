package minio_adapter

import (
	"fmt"
	"os"
)

var MinioConfig = &Config{}

type Config struct {
	Address   string
	AccessKey string
	SecretKey string
	Bucket    string
}

func (c *Config) EnvParse(prefix string) error {
	const (
		configError        = "from Minio config error: %w"
		configEnvMustBeSet = "%s must be set"
	)

	var exists bool
	nameEndpointEnv := prefix + "ADDRESS"
	c.Address, exists = os.LookupEnv(nameEndpointEnv)
	if !exists {
		return fmt.Errorf(configError, fmt.Errorf(configEnvMustBeSet, nameEndpointEnv))
	}

	nameAccessKeyEnv := prefix + "ACCESS_KEY"
	c.AccessKey, exists = os.LookupEnv(nameAccessKeyEnv)
	if !exists {
		return fmt.Errorf(configError, fmt.Errorf(configEnvMustBeSet, nameAccessKeyEnv))
	}

	nameSecretKeyEnv := prefix + "SECRET_KEY"
	c.SecretKey, exists = os.LookupEnv(nameSecretKeyEnv)
	if !exists {
		return fmt.Errorf(configError, fmt.Errorf(configEnvMustBeSet, nameSecretKeyEnv))
	}

	nameBucketEnv := prefix + "BUCKET"
	c.Bucket, exists = os.LookupEnv(nameBucketEnv)
	if !exists {
		return fmt.Errorf(configError, fmt.Errorf(configEnvMustBeSet, nameBucketEnv))
	}

	return nil
}
