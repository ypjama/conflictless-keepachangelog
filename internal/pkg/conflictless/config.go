package conflictless

type bump uint8

const (
	bumpPatch bump = iota
	bumpMinor
	bumpMajor
)

type flagCollection struct {
	bump      *string
	directory *string
	command   string
}

type config struct {
	flags         flagCollection
	bump          bump
	changelogFile string
}
