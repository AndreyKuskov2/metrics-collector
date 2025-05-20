package collector

import (
	"math/rand/v2"
	"runtime"
	"strconv"

	"github.com/AndreyKuskov2/metrics-collector/internal/models"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func toFloat64Pointer(value float64) *float64 {
	return &value
}

func CollectMetrics(pollCount int64) map[string]models.Metrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return map[string]models.Metrics{
		"Alloc": {
			ID:    "Alloc",
			MType: "gauge",
			Value: toFloat64Pointer(float64(m.Alloc)),
		},
		"BuckHashSys": {
			ID:    "BuckHashSys",
			MType: "gauge",
			Value: toFloat64Pointer(float64(m.BuckHashSys)),
		},
		"Frees": {
			ID:    "Frees",
			MType: "gauge",
			Value: toFloat64Pointer(float64(m.Frees)),
		},
		"GCCPUFraction": {
			ID:    "GCCPUFraction",
			MType: "gauge",
			Value: toFloat64Pointer(float64(m.GCCPUFraction)),
		},
		"GCSys": {
			ID:    "GCSys",
			MType: "gauge",
			Value: toFloat64Pointer(float64(m.GCSys)),
		},
		"HeapAlloc": {
			ID:    "HeapAlloc",
			MType: "gauge",
			Value: toFloat64Pointer(float64(m.HeapAlloc)),
		},
		"HeapIdle": {
			ID:    "HeapIdle",
			MType: "gauge",
			Value: toFloat64Pointer(float64(m.HeapIdle)),
		},
		"HeapInuse": {
			ID:    "HeapInuse",
			MType: "gauge",
			Value: toFloat64Pointer(float64(m.HeapInuse)),
		},
		"HeapObjects": {
			ID:    "HeapObjects",
			MType: "gauge",
			Value: toFloat64Pointer(float64(m.HeapObjects)),
		},
		"HeapReleased": {
			ID:    "HeapReleased",
			MType: "gauge",
			Value: toFloat64Pointer(float64(m.HeapReleased)),
		},
		"HeapSys": {
			ID:    "HeapSys",
			MType: "gauge",
			Value: toFloat64Pointer(float64(m.HeapSys)),
		},
		"LastGC": {
			ID:    "LastGC",
			MType: "gauge",
			Value: toFloat64Pointer(float64(m.LastGC)),
		},
		"Lookups": {
			ID:    "Lookups",
			MType: "gauge",
			Value: toFloat64Pointer(float64(m.Lookups)),
		},
		"MCacheInuse": {
			ID:    "MCacheInuse",
			MType: "gauge",
			Value: toFloat64Pointer(float64(m.MCacheInuse)),
		},
		"MCacheSys": {
			ID:    "MCacheSys",
			MType: "gauge",
			Value: toFloat64Pointer(float64(m.MCacheSys)),
		},
		"MSpanInuse": {
			ID:    "MSpanInuse",
			MType: "gauge",
			Value: toFloat64Pointer(float64(m.MSpanInuse)),
		},
		"MSpanSys": {
			ID:    "MSpanSys",
			MType: "gauge",
			Value: toFloat64Pointer(float64(m.MSpanSys)),
		},
		"Mallocs": {
			ID:    "Mallocs",
			MType: "gauge",
			Value: toFloat64Pointer(float64(m.Mallocs)),
		},
		"NextGC": {
			ID:    "NextGC",
			MType: "gauge",
			Value: toFloat64Pointer(float64(m.NextGC)),
		},
		"NumForcedGC": {
			ID:    "NumForcedGC",
			MType: "gauge",
			Value: toFloat64Pointer(float64(m.NumForcedGC)),
		},
		"NumGC": {
			ID:    "NumGC",
			MType: "gauge",
			Value: toFloat64Pointer(float64(m.NumGC)),
		},
		"OtherSys": {
			ID:    "OtherSys",
			MType: "gauge",
			Value: toFloat64Pointer(float64(m.OtherSys)),
		},
		"PauseTotalNs": {
			ID:    "PauseTotalNs",
			MType: "gauge",
			Value: toFloat64Pointer(float64(m.PauseTotalNs)),
		},
		"StackInuse": {
			ID:    "StackInuse",
			MType: "gauge",
			Value: toFloat64Pointer(float64(m.StackInuse)),
		},
		"StackSys": {
			ID:    "StackSys",
			MType: "gauge",
			Value: toFloat64Pointer(float64(m.StackSys)),
		},
		"Sys": {
			ID:    "Sys",
			MType: "gauge",
			Value: toFloat64Pointer(float64(m.Sys)),
		},
		"TotalAlloc": {
			ID:    "TotalAlloc",
			MType: "gauge",
			Value: toFloat64Pointer(float64(m.TotalAlloc)),
		},
		"PollCount": {
			ID:    "PollCount",
			MType: "counter",
			Delta: &pollCount,
		},
		"RandomValue": {
			ID:    "RandomValue",
			MType: "gauge",
			Value: toFloat64Pointer(rand.Float64()),
		},
	}
}

func CollectAdditionMetrics() map[string]models.Metrics {
	metrics := make(map[string]models.Metrics)

	v, _ := mem.VirtualMemory()

	cpuUtilization, _ := cpu.Percent(0, true)

	metrics["TotalMemory"] = models.Metrics{
		ID:    "TotalMemory",
		MType: "gauge",
		Value: toFloat64Pointer(float64(v.Total)),
	}

	metrics["FreeMemory"] = models.Metrics{
		ID:    "FreeMemory",
		MType: "gauge",
		Value: toFloat64Pointer(float64(v.Free)),
	}

	for i, cpuPercent := range cpuUtilization {
		metrics["CPUutilization"+strconv.Itoa(i+1)] = models.Metrics{
			ID:    "CPUutilization" + strconv.Itoa(i+1),
			MType: "gauge",
			Value: toFloat64Pointer(cpuPercent),
		}
	}

	return metrics
}
