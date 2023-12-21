package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/creativeprojects/catalogue/index"
	"github.com/creativeprojects/catalogue/volume"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func init() {
	volumeCmd.AddCommand(volumeAddCmd)
}

var volumeAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new volume to index",
	Long:  "Add new volume to index: please specify a path of the volume as an argument",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			pterm.Error.Println("Please specify the path of the volume to add")
			return
		}

		volumePath := args[0]
		pterm.Info.Printfln("Analyzing volume %q...", volumePath)

		_, err := os.Stat(volumePath)
		if err != nil {
			pterm.Error.Println("Cannot open path specified:", err)
			return
		}

		vol, err := volume.NewVolumeFromPath(volumePath)
		if err != nil {
			pterm.Error.Println("Cannot get volume information:", err)
			return
		}
		volume.PrintVolume(vol)

		indexer := index.NewIndexer(volumePath)
		indexer.Run()
	},
}

func printStat(stat os.FileInfo) {
	fmt.Printf("Name: %s\n", stat.Name())
	fmt.Printf("Type: %s\n", strings.Join(getFileTypes(stat.Mode()), ", "))
	fmt.Printf("Sys: %+v\n", stat.Sys())
}

func getFileTypes(mode os.FileMode) []string {
	fileTypes := make([]string, 0)
	if mode.IsRegular() {
		fileTypes = append(fileTypes, "file")
		// It cannot be any other type
		return fileTypes
	}
	// Might change this to a switch?
	if mode&os.ModeDir != 0 {
		fileTypes = append(fileTypes, "directory")
	}
	if mode&os.ModeNamedPipe != 0 {
		fileTypes = append(fileTypes, "named pipe")
	}
	if mode&os.ModeSocket != 0 {
		fileTypes = append(fileTypes, "socket")
	}
	if mode&os.ModeSymlink != 0 {
		fileTypes = append(fileTypes, "symlink")
	}
	if mode&os.ModeDevice != 0 {
		fileTypes = append(fileTypes, "device")
	}
	if mode&os.ModeCharDevice != 0 {
		fileTypes = append(fileTypes, "character device")
	}
	if mode&os.ModeIrregular != 0 {
		fileTypes = append(fileTypes, "irregular")
	}
	return fileTypes
}
