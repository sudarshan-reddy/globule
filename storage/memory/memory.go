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
	//consider sharding this
	coordList map[int]storage.Coords
}

// New returns a new instance of memoryStore that conforms to
// storage.CoordinateStore which is basically an in memory
// implementation
func New(r io.Reader) (storage.CoordinateStore, error) {
	memStore := memoryStore{
		coordList: make(map[int]storage.Coords),
	}
	reader := csv.NewReader(r)

	//ignore title
	_, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading header: %s", err)
	}

	for {
		records, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error reading csv :%s", err)
		}

		coords, err := makeCoords(records)
		if err != nil {
			return nil, fmt.Errorf("error making coords: %s", err)
		}
		memStore.coordList[coords.ID] = *coords
	}

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

func (m *memoryStore) Get(ID int) (*storage.Coords, error) {
	coords, ok := m.coordList[ID]
	if !ok {
		return nil, storage.ErrRecordNotFound
	}

	return &coords, nil
}

func (m *memoryStore) GetAll() (chan storage.Coords, error) {
	var coordsCh = make(chan storage.Coords)
	go func() {
		defer close(coordsCh)
		for _, coords := range m.coordList {
			coordsCh <- coords
		}
	}()
	return coordsCh, nil
}
