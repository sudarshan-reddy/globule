package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/sudarshan-reddy/globule/configs"
	"github.com/sudarshan-reddy/globule/jobs"
	"github.com/sudarshan-reddy/globule/storage"
	"github.com/sudarshan-reddy/globule/storage/memory"
	"github.com/sudarshan-reddy/globule/storage/psql"
)

func panicOnError(err error, msg string) {
	if err != nil {
		newErr := fmt.Errorf("%s: %s", msg, err)
		panic(newErr)
	}
}

func main() {

	cfg, err := configs.Load()
	panicOnError(err, "error loading configs")

	store, err := getStore(cfg)
	panicOnError(err, "error getting store")

	runner := jobs.New(store,
		storage.Coords{Latitude: cfg.SourceLatitude, Longitude: cfg.SourceLongitude},
		0,
		cfg.MaxParallelJobs,
		cfg.NumberOfPoints)

	nearest, farthest, err := runner.Run()
	panicOnError(err, "error running jobs")

	fmt.Println("the nearest coords are:")
	for _, listing := range nearest {
		fmt.Printf("ID: %d, Distance %f\n", listing.ID, listing.Distance)
	}

	fmt.Println("the furthest coords are:")
	for _, listing := range farthest {
		fmt.Printf("ID: %d, Distance %f\n", listing.ID, listing.Distance)
	}

}

func getStore(cfg *configs.Config) (storage.CoordinateStore, error) {

	switch cfg.InputDataType {
	case "csv":
		f, err := os.Open(cfg.CSVFileName)
		if err != nil {
			return nil, err
		}

		return memory.New(f, cfg.ParallelReads)
	case "psql":
		return psql.New(cfg.DBConnectString)
	default:
		return nil, errors.New("invalid input type: unsupported")
	}

}
