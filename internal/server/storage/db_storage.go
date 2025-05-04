package storage

import (
	"context"
	"database/sql"

	"github.com/AndreyKuskov2/metrics-collector/internal/models"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStorage struct {
	DB *sql.DB
}

func NewDBStorage(cfg *config.ServerConfig) (*DBStorage, error) {
	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	return &DBStorage{
		DB: db,
	}, nil
}

func (s *DBStorage) Ping(ctx context.Context) error {
	if err := s.DB.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (s *DBStorage) GetAllMetrics() (map[string]*models.Metrics, error) {
	return nil, nil
}

func (s *DBStorage) GetMetric(metricName string) (*models.Metrics, bool) {
	// if metric, ok := s.memStorage[metricName]; ok {
	// 	return metric, true
	// }
	return nil, false
}

func (s *DBStorage) UpdateMetric(metric *models.Metrics) error {
	// s.memStorage[metric.ID] = metric
	return nil
}
