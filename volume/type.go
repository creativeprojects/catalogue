package volume

// Type specifies the type of volume
type Type uint8

// Volume type - Don't change the order, it's based on Windows API:
// https://docs.microsoft.com/en-gb/windows/win32/api/fileapi/nf-fileapi-getdrivetypew
//
// Windows types are up until DriveRAM, other types are for unix based systems
const (
	DriveUnknown Type = iota
	DriveInvalid
	DriveRemovable
	DriveFixed
	DriveNetwork
	DriveOptical
	DriveRAM
	DriveLoopback
	DriveVirtual
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
	case DriveVirtual:
		return "Virtual device"
	default:
		return "Unknown"
	}
}
