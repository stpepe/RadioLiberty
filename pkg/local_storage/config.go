package local_storage

import (
	"RadioLiberty/pkg/local_storage/sqlite_adapter"
	"fmt"
	"maps"
	"os"
	"slices"
)

var LocalStorageConfigsRegister = map[string]LocalStorageConfig{
	"sqlite": sqlite_adapter.SqliteConfig,
}

type LocalStorageConfig interface {
	EnvParse(prefix string) error
}

type Config struct {
	StorageType string
}

func (c *Config) EnvParse(prefix string) error {
	const errorMsg = "from local storage EnvParse error: %w"
	const envVariableMustBeSet = "env variable %s must be set"
	const envVariableMustBeSetOneOf = "env variable %s must be set one of: %v"

	var exists bool

	nameStorageTypeEnv := prefix + "TYPE"
	c.StorageType, exists = os.LookupEnv(nameStorageTypeEnv)
	if !exists {
		return fmt.Errorf(errorMsg, fmt.Errorf(envVariableMustBeSet, nameStorageTypeEnv))
	}

	storageConfig, exists := LocalStorageConfigsRegister[c.StorageType]
	if !exists {
		return fmt.Errorf(
			errorMsg,
			fmt.Errorf(
				envVariableMustBeSetOneOf,
				nameStorageTypeEnv,
				slices.Collect(maps.Keys(LocalStorageConfigsRegister)),
			),
		)
	}

	err := storageConfig.EnvParse(prefix)
	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}
	return nil
}
