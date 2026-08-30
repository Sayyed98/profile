package config_test

import (
	"testing"

	"github.com/mohdhujaifa/profile/internal/config"
	"github.com/stretchr/testify/require"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("HTTP_ADDR", "")
	t.Setenv("MYSQL_DSN", "")
	t.Setenv("MYSQL_USER", "")
	t.Setenv("MYSQL_PASSWORD", "")
	t.Setenv("MYSQL_HOST", "")
	t.Setenv("MYSQL_PORT", "")
	t.Setenv("MYSQL_DATABASE", "")
	cfg := config.Load()
	require.Equal(t, ":8080", cfg.HTTPAddr)
	require.Equal(t, 4, cfg.WorkerPoolSize)
	require.Contains(t, cfg.MySQLDSN, "tcp(127.0.0.1:3307)/portfolio")
	require.Equal(t, "root@tcp(127.0.0.1:3307)/portfolio", cfg.MySQLTarget())
}

func TestLoadOverrides(t *testing.T) {
	t.Setenv("HTTP_ADDR", ":9090")
	t.Setenv("WORKER_POOL_SIZE", "8")
	t.Setenv("CORS_ORIGINS", "http://a.com, http://b.com")
	cfg := config.Load()
	require.Equal(t, ":9090", cfg.HTTPAddr)
	require.Equal(t, 8, cfg.WorkerPoolSize)
	require.Equal(t, []string{"http://a.com", "http://b.com"}, cfg.CORSOrigins)
}

func TestMySQLDSNFromParts(t *testing.T) {
	t.Setenv("MYSQL_DSN", "")
	t.Setenv("MYSQL_USER", "portfolio")
	t.Setenv("MYSQL_PASSWORD", "secret")
	t.Setenv("MYSQL_HOST", "127.0.0.1")
	t.Setenv("MYSQL_PORT", "3307")
	t.Setenv("MYSQL_DATABASE", "portfolio")
	cfg := config.Load()
	require.Equal(t, "portfolio:secret@tcp(127.0.0.1:3307)/portfolio?parseTime=true&charset=utf8mb4", cfg.MySQLDSN)
}

func TestMySQLDSNOverride(t *testing.T) {
	t.Setenv("MYSQL_DSN", "root:root@tcp(127.0.0.1:3306)/portfolio?parseTime=true")
	t.Setenv("MYSQL_USER", "root")
	t.Setenv("MYSQL_PORT", "3307")
	cfg := config.Load()
	require.Contains(t, cfg.MySQLDSN, "3306")
	require.Equal(t, "root@tcp(127.0.0.1:3306)/portfolio", cfg.MySQLTarget())
}
