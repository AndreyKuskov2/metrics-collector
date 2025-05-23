package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/caarlos0/env"
	"github.com/spf13/pflag"
)

type ServerConfig struct {
	Address         string `env:"ADDRESS"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	MaxRetries      int
	RetryDelay      time.Duration
	SecretKey       string `env:"KEY"`
}

func NewConfig() (*ServerConfig, error) {
	var serverConfig ServerConfig

	pflag.StringVarP(&serverConfig.Address, "address", "a", "localhost:8080", "server address")
	pflag.IntVarP(&serverConfig.StoreInterval, "store-interval", "i", 300, "time interval in seconds")
	pflag.StringVarP(&serverConfig.FileStoragePath, "file-storage-path", "f", "storage.json", "file storage path")
	pflag.BoolVarP(&serverConfig.Restore, "restore", "r", true, "restore from file")
	pflag.StringVarP(&serverConfig.DatabaseDSN, "database-dsn", "d", "", "database url")
	pflag.StringVarP(&serverConfig.SecretKey, "key", "k", "", "secret key")

	pflag.Parse()

	for _, arg := range pflag.Args() {
		if !strings.HasPrefix(arg, "-") {
			log.Fatalf("Unknown flag: %v", arg)
		}
	}

	if err := env.Parse(&serverConfig); err != nil {
		return nil, fmt.Errorf("failed to get environment variable value")
	}

	serverConfig.MaxRetries = 3
	serverConfig.RetryDelay = 1 * time.Second

	return &serverConfig, nil
}
