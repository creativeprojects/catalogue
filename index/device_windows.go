//go:build windows

package index

import "io/fs"

func deviceID(_ fs.FileInfo) uint64 {
	return 0
}
