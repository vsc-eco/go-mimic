package main

import (
	"context"
	"mimic/modules/aggregate"
	"mimic/modules/api"
	"mimic/modules/db"
	"mimic/modules/db/mimic"
	"mimic/modules/db/mimic/blockdb"
	"mimic/modules/db/mimic/condenserdb"
	"mimic/modules/db/mimic/state"
)

func main() {

	db := db.New(db.NewDbConfig())
	db.Init()
	mimicDb := mimic.New(db)
	mimicDb.Init()

	hiveBlocks := blockdb.New(mimicDb)
	stateDb := state.New(mimicDb)
	condenserDb := condenserdb.New(mimicDb)

	plugins := []aggregate.Plugin{
		hiveBlocks,
		stateDb,
		condenserDb,
	}

	agg := aggregate.New(plugins)

	agg.Init()
	agg.Start().Await(context.Background())
	defer agg.Stop()

	router := api.NewAPIServer()
	router.Init()
	router.Start()
}
