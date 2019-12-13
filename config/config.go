package config

import (
	"flag"
	"fmt"
	"os"
)

const (
	defaultBatchSize  = 500
	defaultConfigPath = "config.json"
)

type Config struct {
	Help            bool
	Version         bool
	Validation      bool
	BatchSize       int
	PartnerFile     string
	LogFile         string
	MigrationType   string
	MigrationEnv    string
	MigrationConfig string
}

func New() (*Config, error) {
	c := new(Config)

	flag.BoolVar(&c.Help, "h", false, "Show help")
	flag.BoolVar(&c.Version, "v", false, "Show version")
	flag.BoolVar(&c.Validation, "validation", false, "Create batches for validation")
	flag.IntVar(&c.BatchSize, "n", 0, "Number of partners per batch")
	flag.StringVar(&c.PartnerFile, "f", "", "File with PartnerIDs")
	flag.StringVar(&c.MigrationType, "t", "", "Migration type")
	flag.StringVar(&c.MigrationEnv, "e", "", "Migration environment")
	flag.StringVar(&c.MigrationConfig, "c", "", "Migration config")

	flag.Parse()

	if c.Help {
		c.PrintHelp()
		os.Exit(0)
	}

	if c.Version {
		c.PrintVersion()
		os.Exit(0)
	}

	if c.PartnerFile == "" {
		return nil, fmt.Errorf("PartnerID source has not been provided. I quit")
	}

	if c.MigrationType == "" {
		return nil, fmt.Errorf("Migration type has not been provided. I quit")
	}

	if c.MigrationEnv == "" {
		return nil, fmt.Errorf("Migration environment has not been provided. I quit")
	}

	if c.BatchSize <= 0 {
		fmt.Println("Number of partners per batch is not valid. I will use default batch size:", defaultBatchSize)
		c.BatchSize = defaultBatchSize
	}

	if c.MigrationConfig == "" {
		fmt.Println("Migration config path is not provided. I will use default batch size:", defaultConfigPath)
	}

	return c, nil
}

func (f *Config) PrintHelp() {
	fmt.Println(`
Usage: splitter [Options]
Options:
	-h Print Help
	-v Print Version
	-n Number of partners per batch
	-f File with PartnerIDs
	-t Migration type
	-e Migration environment`)
}

func (f *Config) PrintVersion() {
	fmt.Println("42")
}
