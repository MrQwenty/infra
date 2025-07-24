package retry

import (
	"math"
	"math/rand"
	"strings"
	"time"
)

const (
	DefaultBaseDelay = 30  // seconds
	MaxDelay        = 3600 // 1 hour max
	JitterFactor    = 0.1  // 10% jitter
)

type ErrorCategory string

const (
	ErrorCategoryRateLimit    ErrorCategory = "rate_limit"
	ErrorCategoryNetworkError ErrorCategory = "network_error"
	ErrorCategoryAPIError     ErrorCategory = "api_error"
	ErrorCategoryInvalidNumber ErrorCategory = "invalid_number"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func CalculateNextRetryDelay(retryCount int, baseDelay int, errorCategory ErrorCategory) int64 {
	if baseDelay <= 0 {
		baseDelay = DefaultBaseDelay
	}
	
	var multiplier float64
	switch errorCategory {
	case ErrorCategoryRateLimit:
		multiplier = math.Pow(2, float64(retryCount)) * 2
	case ErrorCategoryNetworkError:
		multiplier = math.Pow(1.5, float64(retryCount))
	case ErrorCategoryAPIError:
		multiplier = math.Pow(2, float64(retryCount))
	default:
		multiplier = math.Pow(2, float64(retryCount))
	}
	
	delay := float64(baseDelay) * multiplier
	
	jitter := delay * JitterFactor * (rand.Float64()*2 - 1)
	delay += jitter
	
	if delay > MaxDelay {
		delay = MaxDelay
	}
	
	return int64(delay)
}

func CategorizeError(err error) ErrorCategory {
	errStr := strings.ToLower(err.Error())
	if strings.Contains(errStr, "rate limit") || strings.Contains(errStr, "429") {
		return ErrorCategoryRateLimit
	}
	if strings.Contains(errStr, "network") || strings.Contains(errStr, "timeout") {
		return ErrorCategoryNetworkError
	}
	if strings.Contains(errStr, "invalid") && strings.Contains(errStr, "number") {
		return ErrorCategoryInvalidNumber
	}
	return ErrorCategoryAPIError
}

func ShouldRetry(errorCategory ErrorCategory) bool {
	return errorCategory != ErrorCategoryInvalidNumber
}
