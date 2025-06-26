package state

import (
	"mimic/modules/db"
	"mimic/modules/db/mimic"
)

type StateDb struct {
	*db.Collection
}

func New(d *mimic.MimicDb) StateDb {
	return StateDb{db.NewCollection(d.DbInstance, "state")}
}

func (s *StateDb) GetState(key string) (string, error) {
	// Implement the logic to get the state by key
	return "", nil
}

func (s *StateDb) SetState(key string, value string) error {
	// Implement the logic to set the state by key
	return nil
}
