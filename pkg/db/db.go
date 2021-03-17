package db

import (
	"github.com/konrads/go-micros/pkg/model"
)

type DB interface {
	Get(id string) (*model.Star, error)
	SaveAll(stars []model.Star) (int, error)
	Close() error
}
