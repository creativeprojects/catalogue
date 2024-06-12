package cmd

import (
	"os"

	"github.com/creativeprojects/catalogue/volume"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func init() {
	volumeCmd.AddCommand(volumeInfoCmd)
}

var volumeInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "View information on a volume",
	Long:  "Volume information: please specify a path of the volume as an argument",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			pterm.Error.Println("Please specify the path of the volume")
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
	},
}
