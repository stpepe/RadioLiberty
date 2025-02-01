package sqlite_adapter

import (
	"RadioLiberty/pkg/models"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

type SqliteStorage struct {
	*sql.DB
}

func NewSqliteStorage() (*SqliteStorage, error) {
	const errorMsg = "from NewSqliteStorage error: %w"

	db, err := sql.Open("sqlite3", SqliteConfig.FilePath)
	if err != nil {
		return nil, fmt.Errorf(errorMsg, err)
	}

	return &SqliteStorage{db}, nil
}

func (s *SqliteStorage) AddToQueue(audio *models.AudioInfo) error {
	const errorMsg = "from sqlite AddToQueue error: %w"

	_, err := s.Exec(
		"INSERT INTO queue (file_name, audio_name, artist) VALUES (?, ?, ?)",
		audio.FileName,
		audio.AudioName,
		audio.Artist,
	)
	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}
	return nil
}

func (s *SqliteStorage) Next() (*models.AudioInfo, error) {
	const errorMsg = "from sqlite Next error: %w"

	audioInfo := &models.AudioInfo{}
	err := s.QueryRow("SELECT id, file_name, audio_name, artist FROM queue ORDER BY id ASC LIMIT 1").
		Scan(&audioInfo.ID, &audioInfo.FileName, &audioInfo.AudioName, &audioInfo.Artist)
	if err != nil {
		return nil, fmt.Errorf(errorMsg, err)
	}

	_, err = s.Exec("DELETE FROM queue WHERE id = ?", audioInfo.ID)
	if err != nil {
		return nil, fmt.Errorf(errorMsg, err)
	}

	return audioInfo, nil
}
func (s *SqliteStorage) Migrate(pathToMigrations string) error {
	const errorMsg = "from LocalStorageMigrate error: %w"

	if err := goose.SetDialect("sqlite"); err != nil {
		return fmt.Errorf(errorMsg, err)
	}

	if err := goose.Up(s.DB, pathToMigrations); err != nil {
		return fmt.Errorf(errorMsg, err)
	}

	return nil
}
