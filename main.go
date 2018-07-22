package main

import (
	"fmt"
	"os"

	"github.com/sudarshan-reddy/globule/jobs"
	"github.com/sudarshan-reddy/globule/storage"
	"github.com/sudarshan-reddy/globule/storage/memory"
)

const (
	fileName        = "./geoData.csv"
	sourceLatitude  = 51.925146
	sourceLongitude = 4.478617
	maxParallelJobs = 10
	numberOfPoints  = 5
)

func panicOnError(err error, msg string) {
	if err != nil {
		newErr := fmt.Errorf("%s: %s", msg, err)
		panic(newErr)
	}
}

func main() {

	f, err := os.Open(fileName)
	panicOnError(err, "error opening data file")

	mem, err := memory.New(f)
	panicOnError(err, "error creating new memory")

	runner := jobs.New(mem,
		storage.Coords{Latitude: sourceLatitude, Longitude: sourceLongitude},
		0,
		maxParallelJobs,
		numberOfPoints)

	nearest, farthest, err := runner.Run()
	panicOnError(err, "error running jobs")

	fmt.Println("the nearest coords are:")
	fmt.Println(nearest)

	fmt.Println("the farthest coords are:")
	fmt.Println(farthest)

}
