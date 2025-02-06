package k8sprobe

import (
	"sync"
)

// Probe represents a thread-safe structure for tracking validation state and cause.
type Probe struct {
	mu      sync.RWMutex
	isValid bool
	cause   Cause
}

// NewProbe creates a new probe with the specified initial state.
func NewProbe(isValid bool) *Probe {
	return NewProbeWithCause(isValid, "OK")
}

// NewProbeWithCause creates a new probe with the specified initial state.
func NewProbeWithCause(isValid bool, cause Cause) *Probe {
	return &Probe{
		isValid: isValid,
		cause:   cause,
	}
}

// SetValid updates the validation state and cause of the Probe in a thread-safe manner.
func (p *Probe) SetValid(state bool, cause Cause) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.isValid = state
	p.cause = cause
}

// IsValid returns the current validation state and its associated cause in a thread-safe manner.
func (p *Probe) IsValid() (bool, Cause) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.isValid, p.cause
}
