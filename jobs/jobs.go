package jobs

import (
	"sort"
	"sync"

	"github.com/sudarshan-reddy/globule/haversine"
	"github.com/sudarshan-reddy/globule/storage"
)

// Runner is the master orchestrating worker pools to
// parallelize calculation
type Runner struct {
	sync.Mutex
	coordinateStore storage.CoordinateStore
	haversine       *haversine.Haversine
	semaphore       chan struct{}
	sourceCoords    storage.Coords
	coordsList      []CoordsWithDistance
	numberOfPoints  int
}

// CoordsWithDistance ...
type CoordsWithDistance struct {
	storage.Coords
	distance float64
}

// New returns a new instance of Runner
func New(coordinateStore storage.CoordinateStore, sourceCoords storage.Coords,
	radius, maxParallelJobs, numberOfPoints int) *Runner {
	return &Runner{coordinateStore: coordinateStore,
		haversine:      haversine.New(radius),
		semaphore:      make(chan struct{}, maxParallelJobs),
		sourceCoords:   sourceCoords,
		coordsList:     make([]CoordsWithDistance, 0),
		numberOfPoints: numberOfPoints,
	}
}

// Run ...
func (r *Runner) Run() (nearest, farthest []CoordsWithDistance, err error) {
	coordsCh, err := r.coordinateStore.GetAll()
	if err != nil {
		return nil, nil, err
	}
	for coords := range coordsCh {
		r.semaphore <- struct{}{}
		go r.handle(coords)

	}

	sort.Slice(r.coordsList, func(i, j int) bool {
		return r.coordsList[i].distance < r.coordsList[j].distance
	})

	nearest = r.coordsList[:r.numberOfPoints]
	farthest = r.coordsList[len(r.coordsList)-r.numberOfPoints:]
	sort.Slice(farthest, func(i, j int) bool {
		return farthest[i].distance > farthest[j].distance
	})

	return
}

func (r *Runner) handle(coords storage.Coords) {
	distance := r.haversine.ShortestGCD(r.sourceCoords, coords)
	r.coordsList = append(r.coordsList, CoordsWithDistance{Coords: coords, distance: distance})
	<-r.semaphore
}
