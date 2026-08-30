package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/mohdhujaifa/profile/internal/model"
	"github.com/redis/go-redis/v9"
)

const portfolioKey = "portfolio:v1"

type PortfolioCache interface {
	Get(ctx context.Context) (model.Portfolio, bool, error)
	Set(ctx context.Context, p model.Portfolio) error
	Invalidate(ctx context.Context) error
}

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedis(addr, password string, ttl time.Duration) *RedisCache {
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	return &RedisCache{client: c, ttl: ttl}
}

func (c *RedisCache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

func (c *RedisCache) Get(ctx context.Context) (model.Portfolio, bool, error) {
	var p model.Portfolio
	raw, err := c.client.Get(ctx, portfolioKey).Bytes()
	if err == redis.Nil {
		return p, false, nil
	}
	if err != nil {
		return p, false, err
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return p, false, err
	}
	return p, true, nil
}

func (c *RedisCache) Set(ctx context.Context, p model.Portfolio) error {
	raw, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, portfolioKey, raw, c.ttl).Err()
}

func (c *RedisCache) Invalidate(ctx context.Context) error {
	return c.client.Del(ctx, portfolioKey).Err()
}

type MemoryCache struct {
	value *model.Portfolio
}

func NewMemory() *MemoryCache {
	return &MemoryCache{}
}

func (c *MemoryCache) Get(_ context.Context) (model.Portfolio, bool, error) {
	if c.value == nil {
		return model.Portfolio{}, false, nil
	}
	return *c.value, true, nil
}

func (c *MemoryCache) Set(_ context.Context, p model.Portfolio) error {
	cp := p
	c.value = &cp
	return nil
}

func (c *MemoryCache) Invalidate(_ context.Context) error {
	c.value = nil
	return nil
}
