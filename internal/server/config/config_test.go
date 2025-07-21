package config

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
)

func resetFlags() {
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
}

func TestNewConfig_Defaults(t *testing.T) {
	resetFlags()
	os.Clearenv()
	cfg, err := NewConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Address != "localhost:8080" {
		t.Errorf("expected default address, got %s", cfg.Address)
	}
	if cfg.StoreInterval != 300 {
		t.Errorf("expected default store interval 300, got %d", cfg.StoreInterval)
	}
	if cfg.FileStoragePath != "storage.json" {
		t.Errorf("expected default file storage path, got %s", cfg.FileStoragePath)
	}
	if !cfg.Restore {
		t.Errorf("expected default restore true")
	}
	if cfg.DatabaseDSN != "" {
		t.Errorf("expected default database dsn empty, got %s", cfg.DatabaseDSN)
	}
	if cfg.SecretKey != "" {
		t.Errorf("expected default secret key empty, got %s", cfg.SecretKey)
	}
	if cfg.MaxRetries != 3 {
		t.Errorf("expected MaxRetries 3, got %d", cfg.MaxRetries)
	}
}

func TestNewConfig_EnvOverride(t *testing.T) {
	resetFlags()
	os.Clearenv()
	os.Setenv("ADDRESS", "1.2.3.4:9999")
	os.Setenv("STORE_INTERVAL", "42")
	os.Setenv("FILE_STORAGE_PATH", "foo.json")
	os.Setenv("RESTORE", "false")
	os.Setenv("DATABASE_DSN", "pg://test")
	os.Setenv("KEY", "supersecret")
	cfg, err := NewConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Address != "1.2.3.4:9999" {
		t.Errorf("env override failed for address, got %s", cfg.Address)
	}
	if cfg.StoreInterval != 42 {
		t.Errorf("env override failed for store interval, got %d", cfg.StoreInterval)
	}
	if cfg.FileStoragePath != "foo.json" {
		t.Errorf("env override failed for file storage path, got %s", cfg.FileStoragePath)
	}
	if cfg.Restore {
		t.Errorf("env override failed for restore, got true")
	}
	if cfg.DatabaseDSN != "pg://test" {
		t.Errorf("env override failed for database dsn, got %s", cfg.DatabaseDSN)
	}
	if cfg.SecretKey != "supersecret" {
		t.Errorf("env override failed for secret key, got %s", cfg.SecretKey)
	}
}

func TestNewConfig_FlagOverride(t *testing.T) {
	resetFlags()
	os.Clearenv()
	os.Args = []string{"cmd", "--address=9.9.9.9:1234", "--store-interval=77", "--file-storage-path=bar.json", "--restore=false", "--database-dsn=pg://flag", "--key=flagsecret"}
	cfg, err := NewConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Address != "9.9.9.9:1234" {
		t.Errorf("flag override failed for address, got %s", cfg.Address)
	}
	if cfg.StoreInterval != 77 {
		t.Errorf("flag override failed for store interval, got %d", cfg.StoreInterval)
	}
	if cfg.FileStoragePath != "bar.json" {
		t.Errorf("flag override failed for file storage path, got %s", cfg.FileStoragePath)
	}
	if cfg.Restore {
		t.Errorf("flag override failed for restore, got true")
	}
	if cfg.DatabaseDSN != "pg://flag" {
		t.Errorf("flag override failed for database dsn, got %s", cfg.DatabaseDSN)
	}
	if cfg.SecretKey != "flagsecret" {
		t.Errorf("flag override failed for secret key, got %s", cfg.SecretKey)
	}
}
