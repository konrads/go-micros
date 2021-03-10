package db

import (
	"fmt"

	model "github.com/konrads/go-micros/pkg/model"
)

type DB struct {
	ports map[string]model.Port
}

func New() DB {
	return DB{ports: make(map[string]model.Port)}
}

func (db *DB) Get(id string) (*model.Port, error) {
	if val, ok := db.ports[id]; ok {
		return &val, nil
	} else {
		return nil, fmt.Errorf("No port for id: %v", id)
	}
}

func (db *DB) SaveAll(ports []model.Port) (int, error) {
	for _, port := range ports {
		db.ports[port.Id] = port
	}

	return len(ports), nil // fmt.Errorf("Failed to save ports: %v", ports)
}
