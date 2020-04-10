package volume

// Type specifies the type of volume
type Type uint8

// Volume type
const (
	VolumeTypeUnknown Type = iota
	FixedDrive
	RemovableDrive
	NetworkDrive
	FloppyDisk
	Optical
)
