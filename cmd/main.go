package main

import (
	"context"
	"log/slog"
	"mimic/lib/utils"
	"mimic/modules/aggregate"
	"mimic/modules/api"
	"mimic/modules/db"
	"mimic/modules/db/mimic"
	"mimic/modules/db/mimic/blockdb"
	"mimic/modules/db/mimic/condenserdb"
	"os"
)

var mimicDb *mimic.MimicDb

func init() {
	// initialize logging
	level := slog.LevelInfo

	switch utils.EnvOrDefault("LOG_LEVEL", "info") {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	slog.SetDefault(slog.New(handler))

	// initialize database
	db := db.New(db.NewDbConfig())
	db.Init()

	mimicDb = mimic.New(db)
	mimicDb.Init()
}

func main() {
	// hiveBlocks := blockdb.New(mimicDb)
	// stateDb := state.New(mimicDb)
	condenserDb := condenserdb.New(mimicDb)
	blockDb := blockdb.New(mimicDb)

	plugins := []aggregate.Plugin{
		// hiveBlocks,
		// stateDb,
		condenserDb,
		blockDb,
	}

	agg := aggregate.New(plugins)

	agg.Init()
	agg.Start().Await(context.Background())
	defer agg.Stop()

	router := api.NewAPIServer()
	router.Init()
	router.Start()
}
