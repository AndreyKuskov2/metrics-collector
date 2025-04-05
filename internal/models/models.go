package models

import "net/http"

type Metric struct {
	Type  string
	Value interface{}
}

type MetricResposne struct {
	Value interface{}
}

// Render implements render.Renderer.
func (m *MetricResposne) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
