package queue

import (
	"RadioLiberty/internal/queue"
	"RadioLiberty/pkg/local_storage"
	"RadioLiberty/pkg/s3"
	"log/slog"
)

func Run() {
	const errorMsg = "from radio Run error"

	generalConfig := &GeneralConfig{}
	err := generalConfig.EnvParse("")
	if err != nil {
		slog.Error(errorMsg, "error", err)
		return
	}

	s3Storage, err := s3.NewS3Storage(generalConfig.S3StorageConfig)
	if err != nil {
		slog.Error(errorMsg, "error", err)
		return
	}

	localStorage, err := local_storage.NewLocalStorage(generalConfig.LocalStorageConfig)
	if err != nil {
		slog.Error(errorMsg, "error", err)
		return
	}

	queueService := queue.NewQueue(localStorage, s3Storage)

	server := queue.NewServer(queueService)
	err = server.Run(generalConfig.Port)
	if err != nil {
		slog.Error(errorMsg, "error", err)
	}
	slog.Info("Server stopped")
}
