// Пакет models содержит основные структуры данных для работы с метриками.
package models

import (
	"fmt"
	"net/http"
)

// Metric - базовая структура хранения метрик.
type Metric struct {
	Type  string
	Value interface{}
}

// Metrics - структура для хранения метрик.
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// Bind - метод для валидации метрик.
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
	return nil
}

// AllMetrics - структура для хранения всех метрик.
type AllMetrics struct {
	// RuntimeMetrics - метрики, которые собираются во время работы программы.
	RuntimeMetrics map[string]Metrics `json:"runtime_metrics"`
	// AdditionalMetrics - дополнительные метрики, которые могут быть переданы в формате JSON.
	AdditionalMetrics map[string]Metrics `json:"additional_metrics"`
}
