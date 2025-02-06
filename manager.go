package k8sprobe

// Cause is a type alias for string, often used to provide reasons or explanations within validation contexts.
type Cause = string

// EmptyCause represents an empty string for Cause, signifying the absence of a specific reason or explanation.
const EmptyCause Cause = ""

// ValidityChecker represents an interface for objects capable of validating their state and returning a boolean result.
// It is used for health checks or determining the validity of various entities within the system.
type ValidityChecker interface {

	// IsValid validates the current state of an object, returning a boolean indicating validity and a string explaining the reason.
	IsValid() (bool, Cause)
}

// Manager provides functionality to manage and monitor health probes for applications through a registry.
type Manager struct {
	registry *registry
}

// NewManager creates and returns a new instance of Manager with an initialized probe registry.
func NewManager() *Manager {
	return &Manager{
		registry: newProbeRegistry(),
	}
}

// RegisterProbe registers a health probe of the specified type with the provided state retriever in the registry.
func (m *Manager) RegisterProbe(probeTypo ProbeType, probe ValidityChecker) {
	m.registry.registerProbe(probeTypo, probe)
}

// CheckProbe checks the status of the specified health probe type and returns true if the probe passes, otherwise false.
func (m *Manager) CheckProbe(probeType ProbeType) (bool, Cause) {
	probes := m.registry.getProbes(probeType)

	for _, probe := range probes {
		if isValid, cause := probe.IsValid(); !isValid {
			return false, cause
		}
	}

	return true, EmptyCause
}
