package db

import (
	"github.com/konrads/go-micros/pkg/model"
)

type DB interface {
	Get(id string) (*model.Port, error)
	SaveAll(ports []model.Port) (int, error)
	Close() error
}
