package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var settingCmd = &cobra.Command{
	Use:   "setting",
	Short: "Set default output directory for torrent files.",
	Long: `Set default output directory for torrent files.
	Write command: rent setting <absolute/path/to/output_directory>.
	By default it $HOME directory.`,
	Args:    cobra.MinimumNArgs(1),
	Run:     SetOutDir,
	Example: "rent setting /home/Видео/",
}

// SetOutDir set output directory for moving files.
func SetOutDir(cmd *cobra.Command, args []string) {
	outPath := args[0]

	stat, err := os.Stat(outPath)
	if os.IsNotExist(err) || stat == nil {
		cobra.CheckErr(err)
	}

	if !stat.IsDir() {
		outPath = filepath.Dir(outPath)
	}

	if outPath[len(outPath)-1] != '/' {
		outPath += "/"
	}

	// Read in the config file
	setConfigFile()

	viper.Set("out_dir", outPath)

	// Save the changes
	err = viper.WriteConfig()
	cobra.CheckErr(err)

	fmt.Println("Output directory has been successfuly selected  :)")
}

// setConfigFile set config file name for writing output directory.
func setConfigFile() {
	viper.AutomaticEnv()

	homeEnv := viper.GetString("HOME")

	viper.SetConfigFile(fmt.Sprintf("%s/.config/rent/settings.yaml", homeEnv))

	err := viper.ReadInConfig()
	cobra.CheckErr(err)
}
