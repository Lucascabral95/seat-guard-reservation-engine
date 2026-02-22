package config

import (
	"testing"
	"time"
)

func TestGetEnvHelpers(t *testing.T) {
	t.Setenv("X_ENV", "value")
	if got := getEnv("X_ENV", "default"); got != "value" {
		t.Fatalf("expected env value, got %q", got)
	}
	if got := getEnv("MISSING", "default"); got != "default" {
		t.Fatalf("expected default value, got %q", got)
	}

	t.Setenv("INT_OK", "42")
	if got := getEnvIntOrDefault("INT_OK", 1); got != 42 {
		t.Fatalf("expected 42, got %d", got)
	}
	t.Setenv("INT_BAD", "abc")
	if got := getEnvIntOrDefault("INT_BAD", 7); got != 7 {
		t.Fatalf("expected fallback 7, got %d", got)
	}

	t.Setenv("DUR_OK", "3m")
	if got := getEnvDurationOrDefault("DUR_OK", time.Second); got != 3*time.Minute {
		t.Fatalf("expected 3m, got %v", got)
	}
	t.Setenv("DUR_BAD", "bad")
	if got := getEnvDurationOrDefault("DUR_BAD", 5*time.Second); got != 5*time.Second {
		t.Fatalf("expected fallback duration, got %v", got)
	}
}

func TestLoadConfig_UsesNewLifetimeVarWithFallback(t *testing.T) {
	t.Setenv("DB_CONN_MAX_LIFETIME", "7m")
	t.Setenv("DB_CONN_MAX_LIFE_TIME", "9m")
	t.Setenv("DB_CONN_MAX_IDLE_TIME", "2m")
	t.Setenv("DB_MAX_OPEN_CONNS", "30")
	t.Setenv("DB_MAX_IDLE_CONNS", "15")

	cfg := LoadConfig()
	if cfg.DbConnMaxLifeTime != 7*time.Minute {
		t.Fatalf("expected new env key precedence, got %v", cfg.DbConnMaxLifeTime)
	}
	if cfg.DbConnMaxIdleTime != 2*time.Minute || cfg.DbMaxOpenConns != 30 || cfg.DbMaxIdleConns != 15 {
		t.Fatalf("unexpected pool config: %+v", cfg)
	}
}
