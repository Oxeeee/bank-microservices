package repo

import "github.com/redis/go-redis/v9"

type BillingCache interface {
}

type billingCache struct {
	redis *redis.Client
}

func NewBillingCache(redis *redis.Client) BillingCache {
	return &billingCache{
		redis: redis,
	}
}
