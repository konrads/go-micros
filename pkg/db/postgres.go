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

func (db *PostgresDB) Get(id string) (*model.Port, error) {
	rows := db.db.QueryRow("SELECT name, coordinates, city, province, country, alias, regions, timezone, unlocs, code FROM ports WHERE id = $1", id)
	var name sql.NullString
	var coordinates []sql.NullFloat64
	var city sql.NullString
	var province sql.NullString
	var country sql.NullString
	var alias []sql.NullString
	var regions []sql.NullString
	var timezone sql.NullString
	var unlocs []sql.NullString
	var code sql.NullString

	err := rows.Scan(&name, pq.Array(&coordinates), &city, &province, &country, pq.Array(&alias), pq.Array(&regions), &timezone, pq.Array(&unlocs), &code)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Fatalf("Failed to unmarshall postgres row due to %v", err)
		return nil, err
	} else {
		res := model.Port{
			Id:          id,
			Name:        name.String,
			Coordinates: toFloat32Arr(coordinates),
			City:        city.String,
			Province:    province.String,
			Country:     country.String,
			Alias:       toStringArr(alias),
			Regions:     toStringArr(regions),
			Timezone:    timezone.String,
			Unlocs:      toStringArr(unlocs),
			Code:        code.String,
		}
		return &res, nil
	}
}

func (db *PostgresDB) SaveAll(ports []model.Port) (int, error) {
	rowsAffected := 0
	for _, p := range ports {
		res, err := db.db.Exec(
			"INSERT INTO ports (id, name, coordinates, city, province, country, alias, regions, timezone, unlocs, code) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) ON CONFLICT(id) DO NOTHING",
			p.Id, p.Name, pq.Array(p.Coordinates), p.City, p.Province, p.Country, pq.Array(p.Alias), pq.Array(p.Regions), p.Timezone, pq.Array(p.Unlocs), p.Code,
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
