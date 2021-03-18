package db

import (
	"github.com/konrads/go-micros/pkg/model"
)

type MemDB struct {
	stars map[string]model.Star
}

func NewMemDB() *MemDB {
	return &MemDB{stars: make(map[string]model.Star)}
}

func (db *MemDB) Get(id string) (*model.Star, error) {
	val, ok := db.stars[id]
	if ok {
		return &val, nil
	}
	return nil, nil
}

func (db *MemDB) SaveAll(stars []model.Star) (int, error) {
	for _, star := range stars {
		db.stars[star.ID] = star
	}
	return len(stars), nil
}

func (db *MemDB) Close() error {
	return nil
}
