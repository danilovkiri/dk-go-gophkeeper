package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config handles all constants and parameters.
type Config struct {
	ServerAddress  string `json:"server_address" env:"SERVER_ADDRESS"`
	DatabaseDSN    string `json:"database_dsn" env:"DATABASE_DSN"`
	UserKey        string `env:"USER_KEY" env-default:"jds__63h3_7ds"`
	AuthBearerName string `env:"BEARER_KEY" env-default:"token"`
}

// NewDefaultConfiguration initializes a configuration struct.
func NewDefaultConfiguration() *Config {
	var cfg Config
	return &cfg
}

// isFlagPassed checks whether the flag was set in CLI
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func (cfg *Config) assignValues(a, c, d *string) error {
	// priority: flag -> env -> json config -> default flag
	var err error
	if *c != "" {
		err = cleanenv.ReadConfig(*c, cfg)
	} else {
		err = cleanenv.ReadEnv(cfg)
	}
	// return err here to stop code execution
	if err != nil {
		return err
	}
	if isFlagPassed("a") || cfg.ServerAddress == "" {
		cfg.ServerAddress = *a
	}
	if isFlagPassed("d") || cfg.DatabaseDSN == "" {
		cfg.DatabaseDSN = *d
	}
	return nil
}

// Parse parses command line arguments and environment and stores them
func (cfg *Config) Parse() error {
	a := flag.String("a", ":8080", "Server address")
	c := flag.String("c", os.Getenv("CONFIG"), "Configuration file path")
	// DatabaseDSN scheme: "postgres://username:password@localhost:5432/database_name"
	d := flag.String("d", "", "Database DSN")
	flag.Parse()
	err := cfg.assignValues(a, c, d)
	return err
}
