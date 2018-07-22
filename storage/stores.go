package storage

import "errors"

// ErrRecordNotFound is a DoesNotExist state
var ErrRecordNotFound = errors.New("record does not exist")

// Coords is a specification to store the input values
type Coords struct {
	ID        int
	Latitude  float64
	Longitude float64
}

// CoordinateStore is the interface for doing coordinate based operations
type CoordinateStore interface {
	GetAll() (chan Coords, error)
}
