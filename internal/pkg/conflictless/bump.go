package conflictless

// Bump is the bump type.
type Bump uint8

const (
	BumpPatch Bump = iota
	BumpMinor
	BumpMajor
)

// InitialVersion returns the initial version for the bump type.
func (b Bump) InitialVersion() string {
	switch b {
	case BumpPatch:
		return "0.0.1"
	case BumpMinor:
		return "0.1.0"
	case BumpMajor:
		return "1.0.0"
	default:
		return "0.1.0"
	}
}
