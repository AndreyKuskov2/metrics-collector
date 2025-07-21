package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/AndreyKuskov2/metrics-collector/internal/models"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/config"
	"github.com/sirupsen/logrus"
)

// Базовая структура для хранения метрик в файле
type FileMemStorage struct {
	FileStorage *os.File
	Encoder     *json.Encoder
	memStorage  map[string]*models.Metrics
	mu          sync.Mutex
}

// Создание структуры FileMemStorage
func NewFileMemStorage() *FileMemStorage {
	return &FileMemStorage{
		memStorage: make(map[string]*models.Metrics),
	}
}

func (s *FileMemStorage) Ping() error {
	return nil
}

// Функция для старта логики хранения метрик в файле
func StartFileStorageLogic(config *config.ServerConfig, s *FileMemStorage, logger *logrus.Logger) {
	if config.FileStoragePath != "" {
		err := s.OpenFile(config.FileStoragePath)
		if err != nil {
			logger.Errorf("Failed to open file: %v", err)
		}
	} else {
		logger.Info("File storage is not specified")
		return
	}

	if config.Restore {
		err := s.LoadMemStorageFromFile()
		if err != nil {
			logger.Errorf("Failed to restore data from file: %v", err)
		}
	}

	go func() {
		for {
			interval := time.Duration(config.StoreInterval) * time.Second
			// if interval == 0 {
			// 	interval = 100 * time.Microsecond // Установите разумное значение по умолчанию
			// }
			time.Sleep(interval)
			s.SaveMemStorageToFile()
		}
	}()
}

func (s *FileMemStorage) OpenFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	encoder := json.NewEncoder(file)

	s.FileStorage = file
	s.Encoder = encoder

	return nil
}

func (s *FileMemStorage) UpdateBatchMetrics(metrics []models.Metrics) error {
	return nil
}

// Сохранение метрик из памяти в файл с его перезаписью
func (s *FileMemStorage) SaveMemStorageToFile() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Очистка файла
	if err := s.FileStorage.Truncate(0); err != nil {
		log.Fatal(err)
		return fmt.Errorf("failed to truncate file: %w", err)
	}

	// Установка указателя файла в начало
	if _, err := s.FileStorage.Seek(0, 0); err != nil {
		log.Fatal(err)
		return fmt.Errorf("failed to seek file: %w", err)
	}

	if err := s.Encoder.Encode(s.memStorage); err != nil {
		log.Fatal(err)
		return fmt.Errorf("failed to encode metrics: %w", err)
	}

	return nil
}

// Загрузка метрик из файла в память
func (s *FileMemStorage) LoadMemStorageFromFile() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Установка указателя файла в начало
	if _, err := s.FileStorage.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}

	// Создание декодера для чтения данных из файла
	decoder := json.NewDecoder(s.FileStorage)

	// Чтение данных из файла
	var metrics map[string]models.Metrics
	for {
		if err := decoder.Decode(&metrics); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("failed to decode metric: %w", err)
		}
		for metricName, metricData := range metrics {
			s.memStorage[metricName] = &metricData
		}
	}

	return nil
}

// Получение всех метрик из памяти
func (s *FileMemStorage) GetAllMetrics() (map[string]*models.Metrics, error) {
	return s.memStorage, nil
}

// Получение конкретной метрики из памяти
func (s *FileMemStorage) GetMetric(metricName string) (*models.Metrics, bool) {
	if metric, ok := s.memStorage[metricName]; ok {
		return metric, true
	}
	return nil, false
}

// Обновление метрики
func (s *FileMemStorage) UpdateMetric(metric *models.Metrics) error {
	s.memStorage[metric.ID] = metric
	return nil
}
