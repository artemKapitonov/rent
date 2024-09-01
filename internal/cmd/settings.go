package cmd

import (
	"fmt"

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

	// Read in the config file
	setConfigFile()

	viper.Set("out_dir", outPath)

	// Save the changes
	err := viper.WriteConfig()
	cobra.CheckErr(err)
}

func setConfigFile() {
	viper.AutomaticEnv()

	homeEnv := viper.GetString("HOME")

	viper.SetConfigFile(fmt.Sprintf("%s/.config/rent/settings.yaml", homeEnv))

	err := viper.ReadInConfig()
	cobra.CheckErr(err)
}
