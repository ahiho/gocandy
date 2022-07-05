package filterparser

// Container represents a container kind.
//
// A zero value for Container corresponds to ContainerNone.
type Container uint

const (
	ContainerNone Container = iota
	ContainerArray
	ContainerMap
)

// ContainerFromString returns a Container based on a string.
//
// Recognized strings: "none", "map", "array".
func ContainerFromString(str string) (c Container, ok bool) {
	switch str {
	case "map":
		return ContainerMap, true
	case "array":
		return ContainerArray, true
	case "none":
		return ContainerNone, true
	default:
		return ContainerNone, false
	}
}
