package db_test

import (
	"testing"

	"github.com/konrads/go-micros/pkg/db"
	"github.com/konrads/go-micros/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestMemDB(t *testing.T) {
	db := db.NewMemDB()
	db.SaveAll([]model.Port{
		{
			Id:          "id1",
			Name:        "name1",
			Coordinates: []float32{1.1, 2.2},
			City:        "city1",
			Province:    "province1",
			Country:     "country1",
			Alias:       []string{"alias1"},
			Regions:     []string{"region1"},
			Timezone:    "timezone1",
			Unlocs:      []string{"unloc1"},
			Code:        "code1",
		},
		{
			Id:          "id2",
			Name:        "name2",
			Coordinates: []float32{1.1, 2.2},
			City:        "city2",
			Province:    "province2",
			Country:     "country2",
			Alias:       []string{"alias2"},
			Regions:     []string{"region2"},
			Timezone:    "timezone2",
			Unlocs:      []string{"unloc2"},
			Code:        "code2",
		},
	})

	res1, err1 := db.Get("id1")
	assert.Nil(t, err1)
	assert.Equal(t, "name1", res1.Name)

	res2, err2 := db.Get("id2")
	assert.Nil(t, err2)
	assert.Equal(t, "city2", res2.City)

	resBogus, errBogus := db.Get("__BOGUS__")
	assert.Nil(t, errBogus)
	assert.Nil(t, resBogus)
}
