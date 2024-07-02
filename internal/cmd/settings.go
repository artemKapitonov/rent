package cmd

import "github.com/spf13/cobra"

var settingCmd = &cobra.Command{
	Use:   "setting",
	Short: "Set default output directory for torrent files",
	Run:   SetOutDir,
}

func SetOutDir(cmd *cobra.Command, args []string) {

}
