package models

import (
	"fmt"
	"net/http"
)

type Metric struct {
	Type  string
	Value interface{}
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m Metrics) Bind(r *http.Request) error {
	if m.ID == "" {
		return fmt.Errorf("id is required parameter")
	}
	if m.MType == "" {
		return fmt.Errorf("type is required parameter")
	}
	if m.MType != "counter" && m.MType != "gauge" {
		return fmt.Errorf("the type field must be one of the following values: counter, gauge")
	}
	// if m.Delta == nil && m.Value == nil {
	// 	return fmt.Errorf("one of the parameters must be passed: delta or value")
	// }
	return nil
}

type AllMetrics struct {
	RuntimeMetrics    map[string]Metrics `json:"runtime_metrics"`
	AdditionalMetrics map[string]Metrics `json:"additional_metrics"`
}
