package audio_processor

import (
	"RadioLiberty/internal/audio_processor"
	"RadioLiberty/pkg/queue_service_client"
	"RadioLiberty/pkg/s3"
	"context"
	"log/slog"
	"os"
	"os/signal"
)

func Run() {
	const errorMsg = "from audio_processor Run error"

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

	queueClient := queue_service_client.NewQueueServiceClient(generalConfig.QueueClientConfig.Host, generalConfig.QueueClientConfig.Port)

	audioProcessor := audio_processor.NewAudioProcessor(queueClient, s3Storage, generalConfig.DefaultAudioName)

	go audioProcessor.Stream(ctx)

	server := audio_processor.NewServer(audioProcessor)
	err = server.Run(generalConfig.Port)
	if err != nil {
		slog.Error(errorMsg, "error", err)
	}
	slog.Info("Server stopped")
}
