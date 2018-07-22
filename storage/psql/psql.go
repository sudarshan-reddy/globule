package psql

import (
	"errors"

	"github.com/sudarshan-reddy/globule/storage"
)

type psqlStore struct {
	storage.CoordinateStore
}

// New is currently a dummy implementation of a psql representation
// of storage.CoordinateStore
func New() (storage.CoordinateStore, error) {
	return &psqlStore{}, nil
}

func (p *psqlStore) Get(ID int) (*storage.Coords, error) {
	return nil, errors.New("not yet implemented")
}

func (p *psqlStore) GetAll() (chan storage.Coords, error) {
	return nil, errors.New("not yet implemented")
}
