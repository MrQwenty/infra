package retry

import (
	"errors"
	"testing"
)

func TestCalculateNextRetryDelay(t *testing.T) {
	tests := []struct {
		name          string
		retryCount    int
		baseDelay     int
		errorCategory ErrorCategory
		expectMin     int64
		expectMax     int64
	}{
		{"first retry network error", 0, 30, ErrorCategoryNetworkError, 25, 35},
		{"second retry rate limit", 1, 30, ErrorCategoryRateLimit, 110, 130},
		{"max delay cap", 10, 30, ErrorCategoryAPIError, MaxDelay-10, MaxDelay+10},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delay := CalculateNextRetryDelay(tt.retryCount, tt.baseDelay, tt.errorCategory)
			if delay < tt.expectMin || delay > tt.expectMax {
				t.Errorf("CalculateNextRetryDelay() = %v, want between %v and %v", delay, tt.expectMin, tt.expectMax)
			}
		})
	}
}

func TestCategorizeError(t *testing.T) {
	tests := []struct {
		error    error
		expected ErrorCategory
	}{
		{errors.New("rate limit exceeded"), ErrorCategoryRateLimit},
		{errors.New("HTTP 429"), ErrorCategoryRateLimit},
		{errors.New("network timeout"), ErrorCategoryNetworkError},
		{errors.New("invalid phone number"), ErrorCategoryInvalidNumber},
		{errors.New("API error"), ErrorCategoryAPIError},
	}
	
	for _, tt := range tests {
		result := CategorizeError(tt.error)
		if result != tt.expected {
			t.Errorf("CategorizeError(%v) = %v, want %v", tt.error, result, tt.expected)
		}
	}
}

func TestShouldRetry(t *testing.T) {
	tests := []struct {
		category ErrorCategory
		expected bool
	}{
		{ErrorCategoryRateLimit, true},
		{ErrorCategoryNetworkError, true},
		{ErrorCategoryAPIError, true},
		{ErrorCategoryInvalidNumber, false},
	}
	
	for _, tt := range tests {
		result := ShouldRetry(tt.category)
		if result != tt.expected {
			t.Errorf("ShouldRetry(%v) = %v, want %v", tt.category, result, tt.expected)
		}
	}
}
