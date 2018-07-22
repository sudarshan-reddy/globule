package haversine

import (
	"math"

	"github.com/sudarshan-reddy/globule/storage"
)

const (
	radius = 6371000
)

// Haversine holds the primitives to compute haversine distance
type Haversine struct {
	radius float64
}

// New returns a new instance of Haversine to calculate distance
// If radius is 0, it defaults to the radius of the earth in metres
func New(inputRadius int) *Haversine {
	h := Haversine{radius: float64(inputRadius)}
	if h.radius == 0 {
		h.radius = radius
	}
	return &h
}

func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}

// ShortestGCD calculates the shortest Great Circle Distance
// between two coordinates on the surface of the Earth.
func (h *Haversine) ShortestGCD(x, y storage.Coords) float64 {
	latitude1 := degreesToRadians(x.Latitude)
	longitude1 := degreesToRadians(x.Longitude)
	latitude2 := degreesToRadians(y.Latitude)
	longitude2 := degreesToRadians(y.Longitude)

	diffLat := latitude2 - latitude1
	diffLon := longitude2 - longitude1

	a := math.Pow(math.Sin(diffLat/2), 2) + math.Cos(latitude1)*math.Cos(latitude2)*
		math.Pow(math.Sin(diffLon/2), 2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return c * h.radius
}
