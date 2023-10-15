package conflictless

// Bump is the bump type.
type Bump uint8

const (
	BumpPatch Bump = iota
	BumpMinor
	BumpMajor
)

func (b Bump) initialVersion() string {
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
