package memory

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sudarshan-reddy/globule/storage"
)

func Test_GetAll(t *testing.T) {
	var testCases = []struct {
		desc           string
		r              io.Reader
		expectedCoords []storage.Coords
		expectedErr    error
	}{
		{
			desc: "read single line without header",
			r:    bytes.NewReader([]byte("id,lat,lng \n382582,37.1768672,-3.608897")),
			expectedCoords: []storage.Coords{storage.Coords{
				ID:        382582,
				Latitude:  37.1768672,
				Longitude: -3.608897,
			},
			},
			expectedErr: nil,
		},
		{
			desc: "read multiple lines without header",
			r:    bytes.NewReader([]byte("id,lat,lng \n382582,37.1768672,-3.608897 \n133, 12.1, -0.11")),
			expectedCoords: []storage.Coords{storage.Coords{
				ID:        382582,
				Latitude:  37.1768672,
				Longitude: -3.608897,
			},
				storage.Coords{
					ID:        133,
					Latitude:  12.1,
					Longitude: -0.11,
				},
			},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			assert := assert.New(t)
			mem, err := New(testCase.r)
			fmt.Println(testCase.expectedErr, err)
			if testCase.expectedErr != nil {
				if testCase.expectedErr.Error() != err.Error() {
					t.Fatal("errors dont match")
				}
			}
			channels, err := mem.GetAll()
			if err != nil {
				t.Fatal(err)
			}

			actualCoords := make([]storage.Coords, 0)
			for ch := range channels {
				actualCoords = append(actualCoords, ch)
			}

			sort.Slice(testCase.expectedCoords, func(i,
				j int) bool {
				return testCase.expectedCoords[i].ID > testCase.expectedCoords[j].ID
			})

			sort.Slice(actualCoords, func(i,
				j int) bool {
				return actualCoords[i].ID > actualCoords[j].ID
			})

			assert.Equal(testCase.expectedCoords, actualCoords)

		})
	}

}
