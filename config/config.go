package config

import (
	"github.com/apex/log"
	"github.com/caarlos0/env"
)

type Config struct {
	MongoURL    string `env:"MONGO_URL" envDefault:"localhost"`
	MongoDBName string `env:"MONGO_DB_NAME" envDefault:"dic"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"debug"`
	Adress      string `env:"adress" envDefault:"localhost:8091"`
	InitFile    string `env:"INIT_FILE" envDefault:"resource/q1_catalog.csv"`
}

var cfg Config

// Get returns a Config instance
func Get() Config {
	if cfg == (Config{}) {
		MustReadFromEnv()
	}
	return cfg
}

// MustReadFromEnv loads config values from environment vars and terminates
// program execution if it fails for any reason
func MustReadFromEnv() {
	err := env.Parse(&cfg)
	if err != nil {
		log.WithError(err).Fatal("reading config")
	}
}
