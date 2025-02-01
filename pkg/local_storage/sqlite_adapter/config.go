package sqlite_adapter

import (
	"log/slog"
	"os"
)

var SqliteConfig = &Config{}

type Config struct {
	FilePath string
}

func (c *Config) EnvParse(prefix string) error {
	var exists bool

	nameFilePathEnv := prefix + "PATH"
	c.FilePath, exists = os.LookupEnv(nameFilePathEnv)
	if !exists {
		c.FilePath = "./local_db/queue.db"
		slog.Warn("env variable %s is not set. Defaulting to %s", nameFilePathEnv, c.FilePath)
	}

	return nil
}
