package storage

import (
	"database/sql"
	"fmt"

	"github.com/AndreyKuskov2/metrics-collector/internal/models"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/config"

	// "github.com/jackc/pgx"

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

	// conn, err := pgx.Connect(context.Background(), cfg.DatabaseDSN)
	// if err != nil {
	// 	// fmt.Println(err)
	// 	// fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	// 	// os.Exit(1)
	// 	return nil, err
	// }
	// defer conn.Close(context.Background())

	// pool, err := pgxpool.New(context.Background(), cfg.DatabaseDSN)
	// if err != nil {
	// 	return nil, err
	// }

	// db := stdlib.OpenDBFromPool(pool)

	return &DBStorage{
		DB: db,
	}, nil
}

func (s *DBStorage) Ping() error {
	if s.DB == nil {
		return fmt.Errorf("database is not connected")
	}
	return s.DB.Ping()
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
