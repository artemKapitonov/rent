package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var settingCmd = &cobra.Command{
	Use:   "setting",
	Short: "Set default output directory for torrent files.",
	Long: `Set default output directory for torrent files.
	Write command: rent setting <absolute/path/to/output_directory>.
	By default it $HOME directory.`,
	Run:     SetOutDir,
	Example: "rent settings /home/Видео",
}

func SetOutDir(cmd *cobra.Command, args []string) {
	outPath := args[0]

	viper.SetConfigFile("/home/.config/rent/settings.yaml")

	// Read in the config file
	err := viper.ReadInConfig()
	cobra.CheckErr(err)
	// Update the value of out_dir
	viper.Set("out_dir", outPath)

	// Save the changes
	err = viper.WriteConfig()
	cobra.CheckErr(err)
}
