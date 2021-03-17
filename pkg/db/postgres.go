package db

import (
	"database/sql"
	"log"

	"github.com/konrads/go-micros/pkg/model"
	"github.com/lib/pq"
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(connUri *string) *PostgresDB {
	// eg: connUri = "postgres://gomicros:password@localhost/gomicros?sslmode=disable"
	db, err := sql.Open("postgres", *connUri)
	if err != nil {
		log.Fatalf("Failed to connect to postgres on uri: %s", *connUri)
	}
	return &PostgresDB{db: db}
}

func (db *PostgresDB) Get(id string) (*model.Star, error) {
	rows := db.db.QueryRow("SELECT name, alias, constellation, coordinates, distance, apparentMagnitude FROM star WHERE id = $1", id)
	var name sql.NullString
	var alias []sql.NullString
	var constellation sql.NullString
	var coordinates []sql.NullFloat64
	var distance sql.NullFloat64
	var apparentMagnitude sql.NullFloat64

	err := rows.Scan(&name, pq.Array(&alias), &constellation, pq.Array(&coordinates), &distance, &apparentMagnitude)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Fatalf("Failed to unmarshall postgres row due to %v", err)
		return nil, err
	} else {
		res := model.Star{
			Id:                id,
			Name:              name.String,
			Alias:             toStringArr(alias),
			Constellation:     constellation.String,
			Coordinates:       toFloat32Arr(coordinates),
			Distance:          float32(distance.Float64),
			ApparentMagnitude: float32(apparentMagnitude.Float64),
		}
		return &res, nil
	}
}

func (db *PostgresDB) SaveAll(stars []model.Star) (int, error) {
	rowsAffected := 0
	for _, s := range stars {
		res, err := db.db.Exec(
			"INSERT INTO ports (id, name, alias, constellation, coordinates, distance, apparentMagnitude) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT(id) DO NOTHING",
			s.Id, s.Name, pq.Array(s.Alias), s.Constellation, pq.Array(s.Coordinates), s.Distance, s.ApparentMagnitude,
		)
		if err != nil {
			return 0, err
		}
		ra, err := res.RowsAffected()
		if err != nil {
			return 0, err
		}
		rowsAffected += int(ra)
	}
	return rowsAffected, nil
}

func (db *PostgresDB) Close() error {
	return db.db.Close()
}

func toStringArr(arr []sql.NullString) []string {
	res := make([]string, len(arr))
	for i, x := range arr {
		res[i] = x.String
	}
	return res
}

func toFloat32Arr(arr []sql.NullFloat64) []float32 {
	res := make([]float32, len(arr))
	for i, x := range arr {
		res[i] = float32(x.Float64)
	}
	return res
}
