package wiring

import (
	"github.com/mjedari/vgang-project/app/configs"
	"github.com/mjedari/vgang-project/domain/contracts"
	"github.com/mjedari/vgang-project/infra/rate_limiter"
)

var Wiring *Wire

type Wire struct {
	Redis       contracts.IStorage
	RateLimiter *rate_limiter.RateLimiter
	Configs     configs.Configuration
}

func NewWire(redis contracts.IStorage, rateLimiter *rate_limiter.RateLimiter, configs configs.Configuration) *Wire {
	return &Wire{Redis: redis, RateLimiter: rateLimiter, Configs: configs}
}

func (w *Wire) GetRateLimiter() *rate_limiter.RateLimiter {
	return w.RateLimiter
}

func (w *Wire) GetStorage() contracts.IStorage {
	return w.Redis
}

func (w *Wire) GetCollector() configs.Collector {
	return w.Configs.Collector
}
