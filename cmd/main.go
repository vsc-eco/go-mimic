package main

import (
	"context"
	"log/slog"
	"mimic/lib/utils"
	"mimic/modules/aggregate"
	"mimic/modules/api"
	"mimic/modules/db"
	"mimic/modules/db/mimic"
	"mimic/modules/db/mimic/accountdb"
	"mimic/modules/db/mimic/blockdb"
	"mimic/modules/db/mimic/transactiondb"
	"mimic/modules/producers"
	"os"
	"time"
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

	plugins := []aggregate.Plugin{
		// hiveBlocks,
		// stateDb,
		blockdb.New(mimicDb),
		accountdb.New(mimicDb),
		transactiondb.New(mimicDb),
		producers.New(),
	}

	agg := aggregate.New(plugins)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	agg.Init()
	agg.Start().Await(ctx)
	defer agg.Stop()

	router := api.NewAPIServer()
	router.Init()
	router.Start()
}
