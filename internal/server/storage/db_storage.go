package storage

import (
	"database/sql"

	"github.com/AndreyKuskov2/metrics-collector/internal/models"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DbStorage struct {
	DB *sql.DB
}

func NewDbStorage(cfg *config.ServerConfig) (*DbStorage, error) {
	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	return &DbStorage{
		DB: db,
	}, nil
}

func (s *DbStorage) Ping() error {
	if err := s.DB.Ping(); err != nil {
		return err
	}
	return nil
}

func (s *DbStorage) GetAllMetrics() (map[string]*models.Metrics, error) {
	return nil, nil
}

func (s *DbStorage) GetMetric(metricName string) (*models.Metrics, bool) {
	// if metric, ok := s.memStorage[metricName]; ok {
	// 	return metric, true
	// }
	return nil, false
}

func (s *DbStorage) UpdateMetric(metric *models.Metrics) error {
	// s.memStorage[metric.ID] = metric
	return nil
}
