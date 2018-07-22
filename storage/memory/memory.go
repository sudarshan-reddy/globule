package memory

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/sudarshan-reddy/globule/storage"
)

type memoryStore struct {
	coordCh chan storage.Coords
}

// New returns a new instance of memoryStore that conforms to
// storage.CoordinateStore which is basically an in memory
// implementation
func New(r io.Reader) (storage.CoordinateStore, error) {
	memStore := memoryStore{
		coordCh: make(chan storage.Coords),
	}
	reader := csv.NewReader(r)

	//ignore title
	_, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading header: %s", err)
	}

	go func(reader *csv.Reader) {
		defer close(memStore.coordCh)
		for {
			records, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				wrappedErr := fmt.Errorf("error reading csv :%s", err)
				if err != nil {
					panic(wrappedErr)
				}
			}

			coords, err := makeCoords(records)
			if err != nil {
				wrappedErr := fmt.Errorf("error making coords: %s", err)
				panic(wrappedErr)
			}
			memStore.coordCh <- *coords
		}
	}(reader)

	return &memStore, nil
}

func makeCoords(records []string) (*storage.Coords, error) {
	if len(records) != 3 {
		return nil, errors.New("data format invalid")
	}

	id, err := strconv.Atoi(strings.Trim(records[0], " "))
	if err != nil {
		return nil, err
	}

	lat, err := strconv.ParseFloat(strings.Trim(records[1], " "), 64)
	if err != nil {
		return nil, err
	}

	lng, err := strconv.ParseFloat(strings.Trim(records[2], " "), 64)
	if err != nil {
		return nil, err
	}

	return &storage.Coords{
		ID:        id,
		Latitude:  lat,
		Longitude: lng,
	}, nil
}

func (m *memoryStore) GetAll() (chan storage.Coords, error) {
	return m.coordCh, nil
}
