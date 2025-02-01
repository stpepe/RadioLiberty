package cmd

import (
	"RadioLiberty/cmd/migrate"
	"RadioLiberty/cmd/run"
	"flag"
	"log/slog"
	"os"
)

func Run() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	migrateFlag := flag.Bool("m", false, "migrate default file and local DB")
	flag.Parse()

	if *migrateFlag {
		err := migrate.Migrate()
		if err != nil {
			slog.Error("from migrate error", "error", err)
			return
		}
	} else {
		run.Run()
	}
}
