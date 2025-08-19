package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env"
	"github.com/spf13/pflag"
)

// ServerConfig - структура для хранения конфигурации сервера.
type ServerConfig struct {
	Address         string `env:"ADDRESS" json:"address"` // Переменная, задающая адрес сервера
	StoreInterval   int    `env:"STORE_INTERVAL" json:"store_interval"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"store_file"`
	Restore         bool   `env:"RESTORE" json:"restore"`
	DatabaseDSN     string `env:"DATABASE_DSN" json:"database_dsn"`
	MaxRetries      int
	RetryDelay      time.Duration
	SecretKey       string `env:"KEY" json:"key"`
	CryptoKey       string `env:"CRYPTO_KEY" json:"crypto_key"`
	Config          string `env:"CONFIG"`
	TrustedSubnet   string `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
	GRPCAddress     string `env:"GRPC_ADDRESS"`
}

// NewConfig - функция для создания новой конфигурации сервера.
func NewConfig() (*ServerConfig, error) {
	var serverConfig ServerConfig

	pflag.StringVarP(&serverConfig.Address, "address", "a", "localhost:8080", "server address")
	pflag.IntVarP(&serverConfig.StoreInterval, "store-interval", "i", 300, "time interval in seconds")
	pflag.StringVarP(&serverConfig.FileStoragePath, "file-storage-path", "f", "storage.json", "file storage path")
	pflag.BoolVarP(&serverConfig.Restore, "restore", "r", true, "restore from file")
	pflag.StringVarP(&serverConfig.DatabaseDSN, "database-dsn", "d", "", "database url")
	pflag.StringVarP(&serverConfig.SecretKey, "key", "k", "", "secret key")
	pflag.StringVar(&serverConfig.CryptoKey, "crypto-key", "", "crypto key")
	pflag.StringVarP(&serverConfig.Config, "config", "c", "", "config")
	pflag.StringVarP(&serverConfig.TrustedSubnet, "trusted-subnet", "t", "", "trusted subnet")
	pflag.StringVarP(&serverConfig.GRPCAddress, "grpc-address", "g", "localhost:3200", "grpc address")

	pflag.Parse()

	if serverConfig.Config != "" {
		file, err := os.ReadFile(serverConfig.Config)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file")
		}
		if err := json.Unmarshal(file, &serverConfig); err != nil {
			return nil, fmt.Errorf("cannot parse json config to object")
		}
	}

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
