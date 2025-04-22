package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
)

type ServerConfig struct {
	Address         string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
}

var (
	address         string
	storeInterval   int
	fileStoragePath string
	restore         bool
)

func init() {
	pflag.StringVarP(&address, "address", "a", "localhost:8080", "server address")
	pflag.IntVarP(&storeInterval, "store-interval", "i", 300, "time interval in seconds")
	pflag.StringVarP(&fileStoragePath, "file-storage-path", "f", "storage.json", "file storage path")
	pflag.BoolVarP(&restore, "restore", "r", true, "restore from file")

	pflag.Parse()

	for _, arg := range pflag.Args() {
		if !strings.HasPrefix(arg, "-") {
			log.Fatalf("Unknown flag: %v", arg)
		}
	}
}

func NewConfig() (*ServerConfig, error) {
	var err error

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		address = envRunAddr
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		storeInterval, err = strconv.Atoi(envStoreInterval)
		if err != nil {
			return nil, err
		}
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		fileStoragePath = envFileStoragePath
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		restore, err = strconv.ParseBool(envRestore)
		if err != nil {
			return nil, err
		}
	}
	return &ServerConfig{
		Address:         address,
		StoreInterval:   storeInterval,
		FileStoragePath: fileStoragePath,
		Restore:         restore,
	}, nil
}
