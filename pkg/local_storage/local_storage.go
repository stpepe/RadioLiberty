package local_storage

import (
	"RadioLiberty/pkg/local_storage/sqlite_adapter"
	"RadioLiberty/pkg/models"
	"fmt"
	"maps"
	"slices"
)

var localStorageRegister = map[string]func() (LocalStorage, error){
	"sqlite": func() (LocalStorage, error) { return sqlite_adapter.NewSqliteStorage() },
}

type LocalStorage interface {
	AddToQueue(audio *models.AudioInfo) error
	Next() (*models.AudioInfo, error)
	Migrate(pathToMigrations string) error
}

func NewLocalStorage(config *Config) (LocalStorage, error) {
	const errorMsg = "from NewLocalStorage error: %w"
	const envVariableMustBeSetOneOf = "env variable %s must be set one of: %v"

	localStorageConstructor, exists := localStorageRegister[config.StorageType]
	if !exists {
		return nil, fmt.Errorf(
			errorMsg,
			fmt.Errorf(
				envVariableMustBeSetOneOf,
				config.StorageType,
				slices.Collect(maps.Keys(localStorageRegister)),
			),
		)
	}

	localStorage, err := localStorageConstructor()
	if err != nil {
		return nil, fmt.Errorf(errorMsg, err)
	}

	return localStorage, nil
}
