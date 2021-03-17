package db

import (
	model "github.com/konrads/go-micros/pkg/model"
)

// unimplemented!!!
type PostgresDB struct {
	connUri *string
}

func NewPostgresDB(connUri *string) *PostgresDB {
	return &PostgresDB{connUri: connUri}
}

func (db *PostgresDB) Get(id string) (*model.Port, error) {
	return nil, nil
}

func (db *PostgresDB) SaveAll(ports []model.Port) (int, error) {
	return 0, nil
}
