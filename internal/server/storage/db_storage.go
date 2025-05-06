package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/AndreyKuskov2/metrics-collector/internal/models"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/config"

	// _ "github.com/jackc/pgx/v5"
	// _ "github.com/lib/pq"

	// "github.com/jackc/pgx/stdlib"
	// "github.com/jackc/pgx/stdlib"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBStorage struct {
	DB  *pgxpool.Pool
	ctx context.Context
}

func NewDBStorage(cfg *config.ServerConfig, ctx context.Context) (*DBStorage, error) {
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return &DBStorage{
		DB:  pool,
		ctx: ctx,
	}, nil
}

func (s *DBStorage) Ping() error {
	return s.DB.Ping(context.Background())
}

func (s *DBStorage) CreateTables() error {
	_, err := s.DB.Exec(s.ctx, createTables)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}

func (s *DBStorage) GetAllMetrics() (map[string]*models.Metrics, error) {
	rows, err := s.DB.Query(s.ctx, getAllMetrics)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics = make(map[string]*models.Metrics)
	for rows.Next() {
		var metric models.Metrics
		err = rows.Scan(&metric.ID, &metric.MType, &metric.Value, &metric.Delta)
		if err != nil {
			return nil, fmt.Errorf("failed to scan metrics: %w", err)
		}
		metrics[metric.ID] = &metric
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over metrics: %w", err)
	}
	return metrics, nil
}

func (s *DBStorage) GetMetric(metricName string) (*models.Metrics, bool) {
	var metric models.Metrics
	if err := s.DB.QueryRow(s.ctx, getMetricByName, metricName).Scan(&metric.ID, &metric.MType, &metric.Value, &metric.Delta); err == nil {
		return &metric, true
	}

	return nil, false
}

func (s *DBStorage) UpdateMetric(metric *models.Metrics) error {
	_, err := s.DB.Exec(s.ctx, insertMetrics, metric.MType, metric.ID, metric.Value, metric.Delta, time.Now())
	if err != nil {
		return fmt.Errorf("failed to insert metric: %w", err)
	}
	return nil
}
