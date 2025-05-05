package storage

import (
	"context"
	"database/sql"

	"github.com/AndreyKuskov2/metrics-collector/internal/models"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/config"

	// _ "github.com/jackc/pgx/v5"
	// _ "github.com/lib/pq"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStorage struct {
	DB  *sql.DB
	DSN string
}

func NewDBStorage(cfg *config.ServerConfig) (*DBStorage, error) {
	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	// err = db.Ping()
	// if err != nil {
	// 	return nil, err
	// }

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
		DB:  db,
		DSN: cfg.DatabaseDSN,
	}, nil
}

// func (s *DBStorage) CloseDB() error {
// 	if err := s.DB.Close(); err != nil {
// 		return err
// 	}
// 	return nil
// }

func (s *DBStorage) Ping(ctx context.Context) error {
	db, err := sql.Open("pgx", s.DSN)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.PingContext(ctx)
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
