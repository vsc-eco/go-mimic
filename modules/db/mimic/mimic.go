package mimic

import (
	a "mimic/modules/aggregate"
	"mimic/modules/db"
)

type MimicDb struct {
	*db.DbInstance
}

var _ a.Plugin = &MimicDb{}

func New(d db.Db) *MimicDb {
	var dbPath = "go-mimic"

	return &MimicDb{db.NewDbInstance(d, dbPath)}
}
