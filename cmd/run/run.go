package run

import (
	"RadioLiberty/internal/audio_processor"
	audio_stream "RadioLiberty/internal/audio_streamer"
	"RadioLiberty/internal/queue"
	"RadioLiberty/internal/server"
	"RadioLiberty/pkg/local_storage"
	"RadioLiberty/pkg/s3"
	"context"
	"log/slog"
	"os"
	"os/signal"
)

func Run() {
	const errorMsg = "from radio Run error"

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		os.Kill,
	)
	defer cancel()

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

	queueService, err := queue.NewQueue(localStorage, s3Storage)
	if err != nil {
		slog.Error(errorMsg, "error", err)
		return
	}
	go queueService.Run(ctx)

	streamer := audio_stream.NewAudioStreamer(queueService.QueueChannel)
	audioChan := streamer.Run(ctx)

	audioProcessor := audio_processor.NewAudioProcessor(audioChan)

	server := server.NewServer(queueService, audioProcessor)
	err = server.Run(generalConfig.Port)
	if err != nil {
		slog.Error(errorMsg, "error", err)
	}
	slog.Info("Server stopped")
}
