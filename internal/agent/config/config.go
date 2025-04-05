package config

import (
	"log"
	"strings"

	"github.com/spf13/pflag"
)

var (
	address        string
	pollInterval   int
	reportInterval int
)

func init() {
	pflag.StringVarP(&address, "address", "a", "localhost:8080", "server address")
	pflag.IntVarP(&pollInterval, "pollInterval", "p", 2, "poll interval in seconds")
	pflag.IntVarP(&reportInterval, "reportInterval", "r", 10, "report interval in seconds")

	pflag.Parse()

	for _, arg := range pflag.Args() {
        if !strings.HasPrefix(arg, "-") {
            log.Fatalf("Unknown flag: %v", arg)
        }
    }
}

type AgentConfig struct {
	Address        string
	ReportInterval int
	PollInterval   int
}

func NewConfig() *AgentConfig {
	return &AgentConfig{
		Address:        address,
		ReportInterval: reportInterval,
		PollInterval:   pollInterval,
	}
}
