package s3

import (
	"RadioLiberty/pkg/s3/minio_adapter"
	"fmt"
	"maps"
	"os"
	"slices"
)

var S3ConfigsRegister = map[string]S3Config{
	"minio": minio_adapter.MinioConfig,
}

type S3Config interface {
	EnvParse(prefix string) error
}

type Config struct {
	S3StorageType string
}

func (c *Config) EnvParse(prefix string) error {
	const errorMsg = "from s3 storage EnvParse error: %w"
	const envVariableMustBeSet = "env variable %s must be set"
	const envVariableMustBeSetOneOf = "env variable %s must be set one of: %v"

	var exists bool

	nameS3StorageTypeEnv := prefix + "TYPE"
	c.S3StorageType, exists = os.LookupEnv(nameS3StorageTypeEnv)
	if !exists {
		return fmt.Errorf(errorMsg, fmt.Errorf(envVariableMustBeSet, nameS3StorageTypeEnv))
	}

	storageConfig, exists := S3ConfigsRegister[c.S3StorageType]
	if !exists {
		return fmt.Errorf(
			errorMsg,
			fmt.Errorf(
				envVariableMustBeSetOneOf,
				nameS3StorageTypeEnv,
				slices.Collect(maps.Keys(S3ConfigsRegister)),
			),
		)
	}

	err := storageConfig.EnvParse(prefix)
	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}
	return nil
}
