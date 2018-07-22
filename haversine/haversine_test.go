package haversine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sudarshan-reddy/globule/storage"
)

func Test_ShortestGCD(t *testing.T) {
	var testCases = []struct {
		desc         string
		coords1      storage.Coords
		coords2      storage.Coords
		expectedDist float64
	}{
		{
			"Velachery to Madipakkam",
			storage.Coords{Latitude: 12.97, Longitude: 80.22},
			storage.Coords{Latitude: 12.96, Longitude: 80.20},
			2435.8181752839755,
		},
		{
			"Rio to Bangkok",
			storage.Coords{Latitude: 22.55, Longitude: 43.12},
			storage.Coords{Latitude: 13.45, Longitude: 100.28},
			6094544.408786774,
		},
		{
			"Port Louis to Padang",
			storage.Coords{Latitude: 20.10, Longitude: 57.30},
			storage.Coords{Latitude: 0.57, Longitude: 100.21},
			5145525.771394785,
		},
		{
			"Oxford to the Vatican",
			storage.Coords{Latitude: 51.45, Longitude: 1.15},
			storage.Coords{Latitude: 41.54, Longitude: 12.27},
			1389179.311829307,
		},
		{
			"Windhoek to Rotterdam",
			storage.Coords{Latitude: 22.34, Longitude: 17.05},
			storage.Coords{Latitude: 51.56, Longitude: 4.29},
			3429893.10043882,
		},
		{
			"Esperanza to Luanda",
			storage.Coords{Latitude: 63.24, Longitude: 56.59},
			storage.Coords{Latitude: 8.50, Longitude: 13.14},
			6996185.95539861,
		},
		{
			"Poles to Paris",
			storage.Coords{Latitude: 90.00, Longitude: 0.00},
			storage.Coords{Latitude: 48.51, Longitude: 2.21},
			4613477.506482742,
		},
		{
			"Turin to Kuala Lumpur",
			storage.Coords{Latitude: 45.04, Longitude: 7.42},
			storage.Coords{Latitude: 3.09, Longitude: 101.42},
			10078111.954385415,
		},
	}

	h := New(0.0)
	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			assert := assert.New(t)
			dist := h.ShortestGCD(test.coords1, test.coords2)
			assert.Equal(test.expectedDist, dist)
		})
	}
}
