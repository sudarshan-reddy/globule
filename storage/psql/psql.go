package psql

import (
	"database/sql"
	"fmt"

	//pq is the sql driver for database/sql
	_ "github.com/lib/pq"

	"github.com/sudarshan-reddy/globule/storage"
)

type psqlStore struct {
	db *sql.DB
	storage.CoordinateStore
}

// New is currently a dummy implementation of a psql representation
// of storage.CoordinateStore
func New(sqlURL string) (storage.CoordinateStore, error) {
	db, err := sql.Open("postgres", sqlURL)
	if err != nil {
		return nil, err
	}
	return &psqlStore{db: db}, nil
}

func (p *psqlStore) GetAll() (chan storage.Coords, error) {
	query := `
		SELECT  id, latitude, longitude 
		FROM coords_list
	`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	coordsCh := make(chan storage.Coords)
	go func() {
		defer close(coordsCh)
		for rows.Next() {
			var coords storage.Coords
			err = rows.Scan(&coords.ID, &coords.Latitude, &coords.Longitude)
			if err != nil {
				wrappedErr := fmt.Errorf("error scanning: %s", err)
				panic(wrappedErr)
			}
			coordsCh <- coords
		}
	}()

	return coordsCh, nil
}
