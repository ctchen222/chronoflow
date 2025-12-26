package service

import "time"

// TimeProvider defines an interface for getting the current time
// This allows for easy testing by injecting a mock time provider
type TimeProvider interface {
	Now() time.Time
	Today() time.Time
}

// RealTimeProvider implements TimeProvider using the system clock
type RealTimeProvider struct{}

// NewRealTimeProvider creates a new RealTimeProvider
func NewRealTimeProvider() *RealTimeProvider {
	return &RealTimeProvider{}
}

func (p *RealTimeProvider) Now() time.Time {
	return time.Now()
}

func (p *RealTimeProvider) Today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
}

// MockTimeProvider implements TimeProvider with a fixed time for testing
type MockTimeProvider struct {
	fixedTime time.Time
}

// NewMockTimeProvider creates a new MockTimeProvider with the given fixed time
func NewMockTimeProvider(t time.Time) *MockTimeProvider {
	return &MockTimeProvider{fixedTime: t}
}

func (p *MockTimeProvider) Now() time.Time {
	return p.fixedTime
}

func (p *MockTimeProvider) Today() time.Time {
	return time.Date(p.fixedTime.Year(), p.fixedTime.Month(), p.fixedTime.Day(), 0, 0, 0, 0, time.Local)
}

// SetTime allows changing the fixed time (useful in tests)
func (p *MockTimeProvider) SetTime(t time.Time) {
	p.fixedTime = t
}
