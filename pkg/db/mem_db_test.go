package db_test

import (
	"testing"

	"github.com/konrads/go-micros/pkg/db"
	"github.com/konrads/go-micros/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestMemDB(t *testing.T) {
	db := db.NewMemDB()
	db.SaveAll([]model.Star{
		{
			ID: "id1",
			StarReq: &model.StarReq{
				Name:              "name1",
				Alias:             []string{"alias1"},
				Constellation:     "constellation1",
				Coordinates:       []float32{1.1, 1.2},
				Distance:          1.1,
				ApparentMagnitude: 11.11,
			},
		},
		{
			ID: "id2",
			StarReq: &model.StarReq{
				Name:              "name2",
				Alias:             []string{"alias2"},
				Constellation:     "constellation2",
				Coordinates:       []float32{2.1, 2.2},
				Distance:          2.2,
				ApparentMagnitude: 22.22,
			},
		},
	})

	res1, err1 := db.Get("id1")
	assert.Nil(t, err1)
	assert.Equal(t, "name1", res1.Name)

	res2, err2 := db.Get("id2")
	assert.Nil(t, err2)
	assert.Equal(t, "constellation2", res2.Constellation)

	resBogus, errBogus := db.Get("__BOGUS__")
	assert.Nil(t, errBogus)
	assert.Nil(t, resBogus)
}
