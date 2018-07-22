package configs

import "github.com/kelseyhightower/envconfig"

// Config holds all the environment variables
type Config struct {
	MainConfig
	CSVConfig
	PSQLConfig
}

// MainConfig holds all the configs needed for Globule
type MainConfig struct {
	SourceLatitude  float64 `envconfig:"SOURCE_LATITUDE" default:"51.925146"`
	SourceLongitude float64 `envconfig:"SOURCE_LONGITUDE" default:"4.478617"`
	MaxParallelJobs int     `envconfig:"MAX_PARALLEL_JOBS" default:"10"`
	NumberOfPoints  int     `envconfig:"NUMBER_OF_POINTS" default:"5"`
	InputDataType   string  `envconfig:"INPUT_DATA_FROM" default:"csv"`
}

// CSVConfig ...
type CSVConfig struct {
	CSVFileName string `envconfig:"CSV_FILE_NAME" default:"./geoData.csv"`
}

// PSQLConfig ...
type PSQLConfig struct {
	DBConnectString string `envconfig:"PSQL_DB_URL"`
}

// Load loads all the configs needed to run globule
func Load() (*Config, error) {
	var config Config
	err := envconfig.Process("GLOBULE", &config.MainConfig)
	if err != nil {
		return nil, err
	}
	switch config.InputDataType {
	case "csv":
		err := envconfig.Process("GLOBULE", &config.CSVConfig)
		if err != nil {
			return nil, err
		}
	case "psql":
		err := envconfig.Process("GLOBULE", &config.PSQLConfig)
		if err != nil {
			return nil, err
		}

	default:
	}

	return &config, err
}
