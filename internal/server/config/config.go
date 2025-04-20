package config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

type ServerConfig struct {
	Address string
}

var (
	address string
)

func init() {
	pflag.StringVarP(&address, "address", "a", "localhost:8080", "server address")

	pflag.Parse()

	for _, arg := range pflag.Args() {
		if !strings.HasPrefix(arg, "-") {
			log.Fatalf("Unknown flag: %v", arg)
		}
	}
}

func NewConfig() *ServerConfig {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		address = envRunAddr
	}
	return &ServerConfig{
		Address: address,
	}
}
