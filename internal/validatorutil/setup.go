package validatorutil

import (
	"log/slog"
	"sync"

	"github.com/go-playground/validator/v10"
)

type paramCache struct {
	mu sync.RWMutex
	paramData map[string]any
}

func SetupCustomValidations(validator *validator.Validate, logger *slog.Logger) {
	// store for validations which have a custom param to cache
	c := &paramCache{
		paramData: make(map[string]any),
	}

	validator.RegisterValidation(durationValidatorTag, durationValidator)
	validator.RegisterValidation(intervalGranularityValidatorTag, makeIntervalGranularityValidator(c, logger))
}

func getParam[T any](cache *paramCache, validatorTag string, param string, parseFn func(string) (T, bool)) (T, bool) {
	cacheKey := validatorTag + "|" + param

	cache.mu.RLock()
	parsedParam, ok := cache.paramData[cacheKey]
	cache.mu.RUnlock()
	if ok {
		return parsedParam.(T), true
	}

	parsedParam, ok = parseFn(param)
	if !ok {
		return parsedParam.(T), false
	}
	cache.mu.Lock()
	cache.paramData[cacheKey] = parsedParam
	cache.mu.Unlock()

	return parsedParam.(T), true
}

