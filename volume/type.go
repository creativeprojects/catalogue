package volume

// Type specifies the type of volume
type Type uint8

// Volume type
const (
	DriveUnknown Type = iota
	DriveFixed
	DriveRemovable
	DriveNetwork
	DriveFloppy
	DriveOptical
)
