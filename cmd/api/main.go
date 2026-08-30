package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mohdhujaifa/profile/internal/cache"
	"github.com/mohdhujaifa/profile/internal/config"
	"github.com/mohdhujaifa/profile/internal/db"
	"github.com/mohdhujaifa/profile/internal/handler"
	"github.com/mohdhujaifa/profile/internal/queue"
	"github.com/mohdhujaifa/profile/internal/repository"
	"github.com/mohdhujaifa/profile/internal/seed"
	"github.com/mohdhujaifa/profile/internal/service"
	"github.com/mohdhujaifa/profile/internal/worker"
)

func main() {
	// Overload so values in .env win over stale exported shell variables.
	if err := godotenv.Overload(); err != nil {
		log.Printf("no .env loaded (%v); using process env and defaults", err)
	}
	cfg := config.Load()
	log.Printf("mysql target: %s", cfg.MySQLTarget())

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	sqlDB, err := db.Open(cfg.MySQLDSN)
	if err != nil {
		log.Fatalf("mysql: %v (is Docker MySQL published on this port? host mysqld often already owns 3306)", err)
	}
	defer sqlDB.Close()

	if err := db.Migrate(ctx, sqlDB); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	repo := repository.NewPortfolioRepository(sqlDB)
	if err := seed.ResumeIfEmpty(ctx, repo); err != nil {
		log.Fatalf("seed: %v", err)
	}

	var portfolioCache cache.PortfolioCache = cache.NewMemory()
	redisCache := cache.NewRedis(cfg.RedisAddr, cfg.RedisPassword, cfg.CacheTTL)
	if err := redisCache.Ping(ctx); err != nil {
		log.Printf("redis unavailable, using in-memory cache: %v", err)
	} else {
		portfolioCache = redisCache
		log.Printf("redis cache enabled at %s", cfg.RedisAddr)
	}

	var publisher queue.Publisher = queue.NoopPublisher{}
	if rabbit, err := queue.NewRabbit(cfg.RabbitMQURL); err != nil {
		log.Printf("rabbitmq unavailable, contact events are local only: %v", err)
	} else {
		publisher = rabbit
		defer rabbit.Close()
		log.Printf("rabbitmq publisher enabled")
	}

	pool := worker.NewPool(ctx, cfg.WorkerPoolSize)
	defer pool.Stop()

	svc := service.NewPortfolioService(repo, portfolioCache, publisher, pool)
	srv := handler.New(svc)

	engine := gin.Default()
	srv.Register(engine, cfg.CORSOrigins)

	httpSrv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("listening on %s", cfg.HTTPAddr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http: %v", err)
		}
	}()

	<-ctx.Done()
	shutdown, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	_ = httpSrv.Shutdown(shutdown)
}
