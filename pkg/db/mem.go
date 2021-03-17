package db

import (
	"fmt"

	"github.com/konrads/go-micros/pkg/model"
)

type MemDB struct {
	ports map[string]model.Port
}

func NewMemDB() *MemDB {
	return &MemDB{ports: make(map[string]model.Port)}
}

func (db *MemDB) Get(id string) (*model.Port, error) {
	if val, ok := db.ports[id]; ok {
		return &val, nil
	} else {
		return nil, fmt.Errorf("No port for id: %v", id)
	}
}

func (db *MemDB) SaveAll(ports []model.Port) (int, error) {
	for _, port := range ports {
		db.ports[port.Id] = port
	}

	return len(ports), nil
}

func (db *MemDB) Close() error {
	return nil
}
