package main

import (
	"context"
	"mimic/modules/aggregate"
	"mimic/modules/api"
	"mimic/modules/db"
	"mimic/modules/db/mimic"
	"mimic/modules/db/mimic/blocks"
	"mimic/modules/db/mimic/state"
)

func main() {

	db := db.New(db.NewDbConfig())
	mimicDb := mimic.New(db)
	hiveBlocks := blocks.New(mimicDb)
	stateDb := state.New(mimicDb)

	plugins := []aggregate.Plugin{
		db,
		mimicDb,
		hiveBlocks,
		stateDb,
	}

	agg := aggregate.New(plugins)

	agg.Init()
	agg.Start().Await(context.Background())
	defer agg.Stop()

	router := api.NewAPIServer()
	router.Init()
	router.Start()
}
