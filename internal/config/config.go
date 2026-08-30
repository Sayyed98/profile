package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTPAddr       string
	MySQLDSN       string
	MySQLUser      string
	MySQLHost      string
	MySQLPort      string
	MySQLDatabase  string
	RedisAddr      string
	RedisPassword  string
	RabbitMQURL    string
	CORSOrigins    []string
	CacheTTL       time.Duration
	WorkerPoolSize int
}

func Load() Config {
	ttl := envInt("CACHE_TTL_SECONDS", 300)
	pool := envInt("WORKER_POOL_SIZE", 4)
	if pool < 1 {
		pool = 1
	}
	origins := strings.Split(env("CORS_ORIGINS", "http://localhost:5173"), ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	user := env("MYSQL_USER", "root")
	pass := env("MYSQL_PASSWORD", "root")
	host := env("MYSQL_HOST", "127.0.0.1")
	port := env("MYSQL_PORT", "3307")
	name := env("MYSQL_DATABASE", "portfolio")
	dsn := os.Getenv("MYSQL_DSN")
	if strings.TrimSpace(dsn) == "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4", user, pass, host, port, name)
	}

	return Config{
		HTTPAddr:       env("HTTP_ADDR", ":8080"),
		MySQLDSN:       dsn,
		MySQLUser:      user,
		MySQLHost:      host,
		MySQLPort:      port,
		MySQLDatabase:  name,
		RedisAddr:      env("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword:  env("REDIS_PASSWORD", ""),
		RabbitMQURL:    env("RABBITMQ_URL", "amqp://guest:guest@127.0.0.1:5672/"),
		CORSOrigins:    origins,
		CacheTTL:       time.Duration(ttl) * time.Second,
		WorkerPoolSize: pool,
	}
}

func (c Config) MySQLTarget() string {
	host, port := c.MySQLHost, c.MySQLPort
	if dsnHost, dsnPort, ok := parseTCPAddr(c.MySQLDSN); ok {
		host, port = dsnHost, dsnPort
	}
	return fmt.Sprintf("%s@tcp(%s:%s)/%s", c.MySQLUser, host, port, c.MySQLDatabase)
}

func parseTCPAddr(dsn string) (string, string, bool) {
	start := strings.Index(dsn, "@tcp(")
	if start < 0 {
		return "", "", false
	}
	rest := dsn[start+len("@tcp("):]
	end := strings.Index(rest, ")")
	if end < 0 {
		return "", "", false
	}
	host, port, err := net.SplitHostPort(rest[:end])
	if err != nil {
		return "", "", false
	}
	return host, port, true
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}
