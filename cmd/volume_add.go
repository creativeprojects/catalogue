package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/creativeprojects/catalogue/volume"

	"github.com/apex/log"
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
		log.SetLevel(log.DebugLevel)
		if len(args) == 0 {
			log.Error("Please specify the path of the volume to add")
			return
		}

		volumePath := args[0]
		log.WithField("path", volumePath).Info("Analyzing volume...")

		_, err := os.Stat(volumePath)
		if err != nil {
			log.WithError(err).Error("Cannot open path specified")
			return
		}
		// printStat(stat)

		vol, err := volume.NewVolumeFromPath(volumePath)
		if err != nil {
			log.WithError(err).Error("Error getting fs stat")
			return
		}
		fmt.Println("")
		volume.PrintVolume(vol)
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
