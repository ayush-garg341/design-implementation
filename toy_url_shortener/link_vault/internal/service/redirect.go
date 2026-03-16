package service

import (
	"context"
	"fmt"
	"time"

	"github.com/linkvault/internal/events"
	"github.com/linkvault/internal/redis"
	"github.com/linkvault/internal/store"
)

type RedirectService struct {
	store    store.LinkStore
	cache    redis.Cache
	producer events.EventProducer
}

func NewRedirectService(store *store.PostgresLinkStore, cache redis.Cache, producer events.EventProducer) *RedirectService {
	return &RedirectService{store, cache, producer}
}

func (svc *RedirectService) GetRedirectUrl(ctx context.Context, shortcode string) (*string, error) {

	// Check the cache
	url, err := svc.cache.Get(ctx, shortcode)
	if err == nil {
		svc.producer.PublishClick(ctx, shortcode)
		return &url, nil
	}

	// Fallback to db
	longUrl, err := svc.store.GetRedirectUrl(ctx, shortcode)
	if err != nil {
		return nil, err
	}

	// populate the cache
	err = svc.cache.Set(ctx, shortcode, *longUrl, 24*time.Hour)
	if err != nil {
		fmt.Println(err.Error())
	}

	// async event
	svc.producer.PublishClick(ctx, shortcode)

	return longUrl, nil
}
