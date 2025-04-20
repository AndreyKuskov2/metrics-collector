package config

import (
	"log"
	"os"
	"strconv"
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
	var err error

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		address = envRunAddr
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		reportInterval, err = strconv.Atoi(envReportInterval)
		if err != nil {
			log.Fatalf("failed to convert reportInterval value to type int")
		}
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		pollInterval, err = strconv.Atoi(envPollInterval)
		if err != nil {
			log.Fatalf("failed to convert pollInterval value to type int")
		}
	}
	return &AgentConfig{
		Address:        address,
		ReportInterval: reportInterval,
		PollInterval:   pollInterval,
	}
}
