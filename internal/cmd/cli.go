package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "rent",
	Short: "Rent - a lightweight CLI torrent client.",
	Long: `Rent is a lightweight torrent client for fast and convenient downloading of files directly from the terminal.
With Rent, you can easily manage download,
view download time, and control download speed without leaving the command line.

Enjoy fast and efficient file downloading with Rent!`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(0)
	}
}

func init() {
	rootCmd.AddCommand(downlaodCmd)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(settingCmd)
}
