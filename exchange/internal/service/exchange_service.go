package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/Oxeeee/bank-microservices/exchange/internal/config"
	customerrors "github.com/Oxeeee/bank-microservices/exchange/pkg/custom_errors"
)

type ExchangeService interface {
	ExchangeRequest(originalCurrType, convertedCurrType string, value float32) (float32, error)
}

type service struct {
	log *slog.Logger
	cfg *config.Config
}

func NewExchangeService(log *slog.Logger, cfg *config.Config) ExchangeService {
	return &service{
		log: log,
		cfg: cfg,
	}
}

var (
	httpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
		Timeout: 5 * time.Second,
	}

	rateCache     = make(map[string]cachedRate)
	cacheMutex    = sync.RWMutex{}
	cacheDuration = 5 * time.Minute
)

type cachedRate struct {
	rate      float64
	timestamp time.Time
}

// TODO: оптимизировать функцию (предположить что дает overhead и пофиксить)
func (s *service) ExchangeRequest(originalCurrType, convertedCurrType string, originalCurrValue float32) (float32, error) {
	cacheKey := fmt.Sprintf("%s:%s", originalCurrType, convertedCurrType)

	cacheMutex.RLock()
	if cr, ok := rateCache[cacheKey]; ok {
		if time.Since(cr.timestamp) < cacheDuration {
			cacheMutex.RUnlock()
			return originalCurrValue * float32(cr.rate), nil
		}
	}
	cacheMutex.RUnlock()

	// получение доступных валют
	url := fmt.Sprintf("https://api.frankfurter.dev/v1/latest?base=%s&symbols=%s", originalCurrType, convertedCurrType)
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return 0, customerrors.ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result struct {
		Rates map[string]float64 `json:"rates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	rate, ok := result.Rates[convertedCurrType]
	if !ok {
		return 0, fmt.Errorf("rate not found for currency: %s", convertedCurrType)
	}
	
	cacheMutex.Lock()
	rateCache[cacheKey] = cachedRate{
		rate: rate,
		timestamp: time.Now(),
	}
	
	cacheMutex.Unlock()
	
	return  originalCurrValue * float32(rate), nil
}
