package jobs

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sudarshan-reddy/globule/storage"
)

func Test_Run(t *testing.T) {
	var testCases = []struct {
		desc     string
		coords   []storage.Coords
		nearest  []CoordsWithDistance
		farthest []CoordsWithDistance
	}{
		{
			desc: "check if in a list of coords, the closest and farthest ones to the housingAnywhere are put through",
			coords: []storage.Coords{
				storage.Coords{Latitude: 12.97, Longitude: 80.22, ID: 1},  //Chennai
				storage.Coords{Latitude: 22.55, Longitude: 43.12, ID: 2},  //Rio
				storage.Coords{Latitude: 13.45, Longitude: 100.28, ID: 3}, //Bangkok
				storage.Coords{Latitude: 20.10, Longitude: 57.30, ID: 4},  //Port Louis
				storage.Coords{Latitude: 0.57, Longitude: 100.21, ID: 5},  //Padang
				storage.Coords{Latitude: 51.45, Longitude: 1.15, ID: 6},   //Oxford
				storage.Coords{Latitude: 41.54, Longitude: 12.27, ID: 7},  //Vatican
				storage.Coords{Latitude: 22.34, Longitude: 17.05, ID: 8},  //Windhoek
				storage.Coords{Latitude: 51.56, Longitude: 4.29, ID: 9},   //Rotterdam
				storage.Coords{Latitude: 63.24, Longitude: 56.59, ID: 10}, //Esperanza
				storage.Coords{Latitude: 8.50, Longitude: 13.14, ID: 11},  //Luanda
				storage.Coords{Latitude: 90.00, Longitude: 0.00, ID: 12},  //North Pole
				storage.Coords{Latitude: 48.51, Longitude: 2.21, ID: 13},  //Paris
				storage.Coords{Latitude: 45.04, Longitude: 7.42, ID: 14},  //Turin
				storage.Coords{Latitude: 3.09, Longitude: 101.42, ID: 15}, //Kuala Lumpur
			},
			nearest: []CoordsWithDistance{
				CoordsWithDistance{Coords: storage.Coords{Latitude: 51.56,
					Longitude: 4.29, ID: 9}, Distance: 42628.635000264374}, //Rotterdam
				CoordsWithDistance{Coords: storage.Coords{Latitude: 51.45,
					Longitude: 1.15, ID: 6}, Distance: 235439.8745583701}, //Oxford
				CoordsWithDistance{Coords: storage.Coords{Latitude: 48.51,
					Longitude: 2.21, ID: 13}, Distance: 412576.0719002216}, //Paris
				CoordsWithDistance{Coords: storage.Coords{Latitude: 45.04,
					Longitude: 7.42, ID: 14}, Distance: 795521.4050696603}, //Turin
				CoordsWithDistance{Coords: storage.Coords{Latitude: 41.54,
					Longitude: 12.27, ID: 7}, Distance: 1.2967591730897105e+06}, //Vatican
			},
			farthest: []CoordsWithDistance{
				CoordsWithDistance{Coords: storage.Coords{Latitude: 0.57,
					Longitude: 100.21, ID: 5}, Distance: 1.035015809436041e+07}, //Padang
				CoordsWithDistance{Coords: storage.Coords{Latitude: 3.09,
					Longitude: 101.42, ID: 15}, Distance: 1.0211367939237328e+07}, //Kuala Lumpur
				CoordsWithDistance{Coords: storage.Coords{Latitude: 13.45,
					Longitude: 100.28, ID: 3}, Distance: 9.225285580451041e+06}, //Bangkok
				CoordsWithDistance{Coords: storage.Coords{Latitude: 12.97,
					Longitude: 80.22, ID: 1}, Distance: 7.9007132932833405e+06}, //Chennai
				CoordsWithDistance{Coords: storage.Coords{Latitude: 20.10,
					Longitude: 57.30, ID: 4}, Distance: 5.742918808309662e+06}, //Port Louis
			},
		},
	}

	for _, testCase := range testCases {
		fs := &fakeStore{fakeValues: testCase.coords}
		runner := New(fs, storage.Coords{Latitude: 51.925146, Longitude: 4.478617},
			0, 5, 5)
		t.Run(testCase.desc, func(t *testing.T) {
			assert := assert.New(t)
			nearest, farthest, err := runner.Run()
			assert.Nil(err)
			assert.Equal(testCase.nearest, nearest)
			assert.Equal(testCase.farthest, farthest)

		})
	}
}

type fakeStore struct {
	storage.CoordinateStore
	fakeValues []storage.Coords
}

func (f *fakeStore) GetAll() (chan storage.Coords, error) {
	coordCh := make(chan storage.Coords)
	go func() {
		defer close(coordCh)
		for _, fakeValue := range f.fakeValues {
			coordCh <- fakeValue
		}
	}()
	return coordCh, nil
}

func benchmarkRun(b *testing.B, numberOfRecords, maxParallelJobs int) {

	for n := 0; n < b.N; n++ {
		fs := &randomFakeStore{numberOfFakeRecords: numberOfRecords}

		runner := New(fs, storage.Coords{Latitude: 0, Longitude: 0}, 0,
			maxParallelJobs, 5)

		runner.Run()
	}
}

type randomFakeStore struct {
	storage.CoordinateStore
	numberOfFakeRecords int
}

func (r *randomFakeStore) GetAll() (chan storage.Coords, error) {
	coordCh := make(chan storage.Coords)
	go func() {
		defer close(coordCh)
		for i := 0; i <= r.numberOfFakeRecords; i++ {
			latitude := rand.Float64()
			longitude := rand.Float64()
			id := rand.Intn(10000)
			coordCh <- storage.Coords{
				ID:        id,
				Latitude:  latitude,
				Longitude: longitude,
			}
		}
	}()
	return coordCh, nil
}

func BenchmarkRunThousandRecordsWith1Goroutine(b *testing.B)   { benchmarkRun(b, 1000, 1) }
func BenchmarkRunThousandRecordsWith10Goroutines(b *testing.B) { benchmarkRun(b, 1000, 10) }
