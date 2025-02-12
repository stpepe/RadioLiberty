package cmd

import (
	"RadioLiberty/cmd/audio_processor"
	"RadioLiberty/cmd/migrate"
	"RadioLiberty/cmd/queue"
	"log/slog"
	"os"
)

func Run() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	switch os.Args[1] {
	case "migrate":
		err := migrate.Migrate()
		if err != nil {
			slog.Error("from migrate error", "error", err)
			return
		}
	case "queue":
		queue.Run()
	case "audio_processor":
		audio_processor.Run()
	default:
		slog.Error("unknown command", "command", os.Args[1])
	}
}
