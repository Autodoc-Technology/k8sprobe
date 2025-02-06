package k8sprobe

import (
	"slices"
	"sync"
)

// registry manages the registration and state retrieval of health probes for application monitoring.
type registry struct {
	mu     sync.RWMutex
	probes map[ProbeType][]ValidityChecker
}

// newProbeRegistry creates a new registry instance.
func newProbeRegistry() *registry {
	return &registry{
		probes: make(map[ProbeType][]ValidityChecker),
	}
}

// registerProbe registers a health probe of a specified type, associating it with the provided ValidityChecker implementation.
func (pr *registry) registerProbe(probeType ProbeType, probe ValidityChecker) {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	pr.probes[probeType] = append(pr.probes[probeType], probe)
}

// getProbes returns all health probes of the specified type.
func (pr *registry) getProbes(probeType ProbeType) []ValidityChecker {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	return slices.Clone(pr.probes[probeType])
}
