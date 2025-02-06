package k8sprobe

// ProbeType defines an integer-based enumeration representing different types of health probes for applications.
type ProbeType int

const (

	// LivenessProbe indicates a probe used to determine if an application is still alive and functioning properly.
	LivenessProbe ProbeType = iota

	// ReadinessProbe indicates a probe used to determine if an application is ready to serve requests.
	ReadinessProbe
)

// String converts a ProbeType to a string.
func (p ProbeType) String() string {
	return [...]string{"LivenessProbe", "ReadinessProbe"}[p]
}
