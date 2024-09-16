package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var settingCmd = &cobra.Command{
	Use:   "setting",
	Short: "Set default output directory for torrent files.",
	Long: `Set default output directory for torrent files.
Write command: rent setting </absolute/path/to/output_directory/>.
By default it is $HOME directory`,
	Args:    cobra.MinimumNArgs(1),
	Run:     SetOutDir,
	Example: "rent setting /home/Video/",
}

// SetOutDir sets the output directory for moving files.
func SetOutDir(_ *cobra.Command, args []string) {
	outPath := args[0]

	if err := validateOutputPath(outPath); err != nil {
		cobra.CheckErr(err)
	}

	// Ensure the output path ends with a '/'
	if !strings.HasSuffix(outPath, "/") {
		outPath += "/"
	}

	setConfigFile()
	viper.Set("out_dir", outPath)

	if err := viper.WriteConfig(); err != nil {
		cobra.CheckErr(err)
	}

	fmt.Fprintln(stderr, "\nOutput directory has been successfully selected :)")
}

// validateOutputPath checks if the provided path exists and is a directory.
func validateOutputPath(path string) error {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", path)
	}

	if !filepath.IsAbs(path) {
		return fmt.Errorf("path is not absolute: %s", path)
	}

	if !stat.IsDir() {
		return fmt.Errorf("path is not a directory: %s", path)
	}

	return nil
}

// setConfigFile sets the config file name for writing the output directory.
func setConfigFile() {
	viper.AutomaticEnv()
	homeEnv := viper.GetString("HOME")
	viper.SetConfigFile(fmt.Sprintf("%s/.config/rent/setting.yaml", homeEnv))

	if err := viper.ReadInConfig(); err != nil {
		cobra.CheckErr(err)
	}
}
