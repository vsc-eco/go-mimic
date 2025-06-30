package main

import (
	"context"
	"mimic/modules/aggregate"
	"mimic/modules/api"
	"mimic/modules/db"
	"mimic/modules/db/mimic"
	"mimic/modules/db/mimic/blockdb"
	"mimic/modules/db/mimic/condenserdb"
)

func main() {

	db := db.New(db.NewDbConfig())
	db.Init()
	mimicDb := mimic.New(db)
	mimicDb.Init()

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
