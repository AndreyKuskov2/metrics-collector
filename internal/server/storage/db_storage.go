package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/AndreyKuskov2/metrics-collector/internal/models"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/config"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBStorage struct {
	DB  *pgxpool.Pool
	ctx context.Context
}

func NewDBStorage(cfg *config.ServerConfig, ctx context.Context) (*DBStorage, error) {
	var pool *pgxpool.Pool
	var err error

	for i := 0; i < cfg.MaxRetries; i++ {
		pool, err = pgxpool.New(context.Background(), cfg.DatabaseDSN)
		if err == nil {
			break
		}
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

func (s *DBStorage) UpdateBatchMetrics(metrics []models.Metrics) error {
	tx, err := s.DB.Begin(s.ctx)
	if err != nil {
		return fmt.Errorf("cannot create transaction")
	}

	counterMap := make(map[string]int64)
	var gauges []models.Metrics

	for _, v := range metrics {
		switch v.MType {
		case utils.COUNTER:
			if v.Delta != nil {
				counterMap[v.ID] += *v.Delta
			}
		case utils.GAUGE:
			gauges = append(gauges, v)
		}
	}

	for k, v := range counterMap {
		_, err = tx.Exec(s.ctx, insertCounterMetrics, utils.COUNTER, k, v, time.Now())
		if err != nil {
			err := tx.Rollback(s.ctx)
			if err != nil {
				return err
			}
			return err
		}
	}

	for _, v := range gauges {
		_, err = tx.Exec(s.ctx, insertGaugeMetrics, utils.GAUGE, v.ID, v.Value, time.Now())
		if err != nil {
			err := tx.Rollback(s.ctx)
			if err != nil {
				return err
			}
			return err
		}
	}

	err = tx.Commit(s.ctx)
	if err != nil {
		return err
	}

	return nil
}
