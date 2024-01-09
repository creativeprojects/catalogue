package volume

// Type specifies the type of volume
type Type uint8

// Volume type - Don't change the order, it's based on Windows API:
// https://docs.microsoft.com/en-gb/windows/win32/api/fileapi/nf-fileapi-getdrivetypew
const (
	DriveUnknown Type = iota
	DriveInvalid
	DriveRemovable
	DriveFixed
	DriveNetwork
	DriveOptical
	DriveRAM
	DriveLoopback
)

func (t Type) String() string {
	switch t {
	case DriveRemovable:
		return "Removable (Floppy, USB drive, Flash card reader)"
	case DriveFixed:
		return "Internal drive"
	case DriveNetwork:
		return "Remote (network) drive"
	case DriveOptical:
		return "Optical drive"
	case DriveRAM:
		return "RAM disk"
	case DriveLoopback:
		return "Loopback (disk image)"
	default:
		return "Unknown"
	}
}
