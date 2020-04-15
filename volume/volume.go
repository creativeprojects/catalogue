package volume

import (
	"fmt"
	"os"
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
	Format          string
	Created         time.Time
	Catalogued      time.Time
	BytesTotal      uint64
	BytesFree       uint64
	RegularFiles    uint64
	HiddenFiles     uint64
	Device          string
	Path            string
	Hostname        string
	IncludeInSearch bool
	Location        string
}

// NewVolumeFromPath creates a populates Volume data from volumePath
func NewVolumeFromPath(volumePath string) (*Volume, error) {
	var err error

	volume := &Volume{
		Created:         time.Now(),
		IncludeInSearch: true,
	}
	hostname, err := os.Hostname()
	volume.Hostname = hostname

	err = getDiskSpace(volumePath, volume)
	if err != nil {
		return nil, err
	}

	err = getFilesystemInfo(volumePath, volume)
	if err != nil {
		return nil, err
	}

	return volume, nil
}

// PrintVolume prints volume information to the console
func PrintVolume(volume *Volume) {
	fmt.Printf("   Hostname: %s\n", volume.Hostname)
	fmt.Printf("    Created: %s\n", volume.Created)
	fmt.Printf("       Name: %s\n", volume.Name)
	fmt.Printf("         ID: %s\n", volume.VolumeID)
	fmt.Printf("     Device: %s\n", volume.Device)
	fmt.Printf("       Type: %s\n", volume.VolumeType.String())
	fmt.Printf("       Path: %s\n", volume.Path)
	fmt.Printf("     Format: %s\n", volume.Format)
	fmt.Printf("Total space: %s\n", getBinaryBytes(volume.BytesTotal))
	fmt.Printf(" Free space: %s\n", getBinaryBytes(volume.BytesFree))
}

func getBinaryBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}
