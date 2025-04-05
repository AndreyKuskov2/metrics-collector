package collector

import (
	"math/rand/v2"
	"runtime"

	"github.com/AndreyKuskov2/metrics-collector/internal/models"
)

func CollectMetrics(pollCount int64) map[string]models.Metric {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return map[string]models.Metric{
		"Alloc": {
			Type:  "gauge",
			Value: float64(m.Alloc),
		},
		"BuckHashSys": {
			Type:  "gauge",
			Value: float64(m.BuckHashSys),
		},
		"Frees": {
			Type:  "gauge",
			Value: float64(m.Frees),
		},
		"GCCPUFraction": {
			Type:  "gauge",
			Value: float64(m.GCCPUFraction),
		},
		"GCSys": {
			Type:  "gauge",
			Value: float64(m.GCSys),
		},
		"HeapAlloc": {
			Type:  "gauge",
			Value: float64(m.HeapAlloc),
		},
		"HeapIdle": {
			Type:  "gauge",
			Value: float64(m.HeapIdle),
		},
		"HeapInuse": {
			Type:  "gauge",
			Value: float64(m.HeapInuse),
		},
		"HeapObjects": {
			Type:  "gauge",
			Value: float64(m.HeapObjects),
		},
		"HeapReleased": {
			Type:  "gauge",
			Value: float64(m.HeapReleased),
		},
		"HeapSys": {
			Type:  "gauge",
			Value: float64(m.HeapSys),
		},
		"LastGC": {
			Type:  "gauge",
			Value: float64(m.LastGC),
		},
		"Lookups": {
			Type:  "gauge",
			Value: float64(m.Lookups),
		},
		"MCacheInuse": {
			Type:  "gauge",
			Value: float64(m.MCacheInuse),
		},
		"MCacheSys": {
			Type:  "gauge",
			Value: float64(m.MCacheSys),
		},
		"MSpanInuse": {
			Type:  "gauge",
			Value: float64(m.MSpanInuse),
		},
		"MSpanSys": {
			Type:  "gauge",
			Value: float64(m.MSpanSys),
		},
		"Mallocs": {
			Type:  "gauge",
			Value: float64(m.Mallocs),
		},
		"NextGC": {
			Type:  "gauge",
			Value: float64(m.NextGC),
		},
		"NumForcedGC": {
			Type:  "gauge",
			Value: float64(m.NumForcedGC),
		},
		"NumGC": {
			Type:  "gauge",
			Value: float64(m.NumGC),
		},
		"OtherSys": {
			Type:  "gauge",
			Value: float64(m.OtherSys),
		},
		"PauseTotalNs": {
			Type:  "gauge",
			Value: float64(m.PauseTotalNs),
		},
		"StackInuse": {
			Type:  "gauge",
			Value: float64(m.StackInuse),
		},
		"StackSys": {
			Type:  "gauge",
			Value: float64(m.StackSys),
		},
		"Sys": {
			Type:  "gauge",
			Value: float64(m.Sys),
		},
		"TotalAlloc": {
			Type:  "gauge",
			Value: float64(m.TotalAlloc),
		},
		"PollCount": {
			Type:  "counter",
			Value: pollCount,
		},
		"RandomValue": {
			Type:  "gauge",
			Value: rand.Float64(),
		},
	}
}
