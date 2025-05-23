package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/caarlos0/env"
	"github.com/spf13/pflag"
)

type AgentConfig struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	MaxRetries     int
	RetryDelay     time.Duration
	SecretKey      string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

func NewConfig() (*AgentConfig, error) {
	var agentConfig AgentConfig

	pflag.StringVarP(&agentConfig.Address, "address", "a", "localhost:8080", "server address")
	pflag.IntVarP(&agentConfig.PollInterval, "pollInterval", "p", 2, "poll interval in seconds")
	pflag.IntVarP(&agentConfig.ReportInterval, "reportInterval", "r", 10, "report interval in seconds")
	pflag.StringVarP(&agentConfig.SecretKey, "key", "k", "", "secret key")
	pflag.IntVarP(&agentConfig.RateLimit, "rateLimit", "l", 0, "rate limit")

	pflag.Parse()

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
