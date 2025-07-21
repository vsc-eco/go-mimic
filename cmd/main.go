package main

import (
	"context"
	"log"
	"log/slog"
	"mimic/lib/utils"
	"mimic/modules/admin"
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

	"github.com/joho/godotenv"
)

var mimicDb *mimic.MimicDb

const (
	mimicServerPort uint16 = 3000
	adminServerPort uint16 = 3001
)

func init() {
	godotenv.Load()

	// initialize logging
	level := slog.LevelInfo

	switch utils.EnvOrDefault("LOG_LEVEL", "info") {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
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
		blockdb.New(mimicDb.Database),
		accountdb.New(mimicDb.Database),
		transactiondb.New(mimicDb.Database),
	}

	dbPlugins := aggregate.New(plugins)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	dbPlugins.Init()
	_, err := dbPlugins.Start().Await(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPlugins.Stop()

	routers := aggregate.New([]aggregate.Plugin{
		api.NewAPIServer(mimicServerPort),
		admin.NewAPIServer(adminServerPort),
		producers.New(),
	})
	routers.Init()
	routers.Start()
	defer routers.Stop()

	select {}
}
