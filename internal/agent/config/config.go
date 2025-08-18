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

// AgentConfig - структура для хранения конфигурации агента.
type AgentConfig struct {
	Address        string `env:"ADDRESS" json:"address"` // Переменная, задающая адрес сервера
	ReportInterval int    `env:"REPORT_INTERVAL" json:"report_interval"`
	PollInterval   int    `env:"POLL_INTERVAL" json:"poll_interval"`
	MaxRetries     int
	RetryDelay     time.Duration
	SecretKey      string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
	CryptoKey      string `env:"CRYPTO_KEY" json:"crypto_key"`
	Config         string `env:"CONFIG"`
}

// NewConfig - функция для создания новой конфигурации агента.
func NewConfig() (*AgentConfig, error) {
	var agentConfig AgentConfig

	pflag.StringVarP(&agentConfig.Address, "address", "a", "localhost:8080", "server address")
	pflag.IntVarP(&agentConfig.PollInterval, "pollInterval", "p", 2, "poll interval in seconds")
	pflag.IntVarP(&agentConfig.ReportInterval, "reportInterval", "r", 10, "report interval in seconds")
	pflag.StringVarP(&agentConfig.SecretKey, "key", "k", "", "secret key")
	pflag.IntVarP(&agentConfig.RateLimit, "rateLimit", "l", 0, "rate limit")
	pflag.StringVar(&agentConfig.CryptoKey, "crypto-key", "", "crypto key")
	pflag.StringVarP(&agentConfig.Config, "config", "c", "", "config")

	pflag.Parse()

	if agentConfig.Config != "" {
		file, err := os.ReadFile(agentConfig.Config)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file")
		}
		if err := json.Unmarshal(file, &agentConfig); err != nil {
			return nil, fmt.Errorf("cannot parse json config to object")
		}
	}

	for _, arg := range pflag.Args() {
		if !strings.HasPrefix(arg, "-") {
			log.Fatalf("Unknown flag: %v", arg)
		}
	}

	if err := env.Parse(&agentConfig); err != nil {
		return nil, fmt.Errorf("failed to get environment variable value")
	}

	agentConfig.MaxRetries = 3
	agentConfig.RetryDelay = 1 * time.Second

	return &agentConfig, nil
}
