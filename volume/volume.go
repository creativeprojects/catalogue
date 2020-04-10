package volume

import (
	"time"

	"github.com/google/uuid"
)

// Volume represents a volume entity
type Volume struct {
	ID              int
	UID             uuid.UUID
	Name            string
	VolumeType      Type
	VolumeID        string
	Created         time.Time
	Catalogued      time.Time
	TotalSize       uint64
	Free            uint64
	RegularFiles    uint64
	HiddenFiles     uint64
	Path            string
	ComputerName    string
	IncludeInSearch bool
}
