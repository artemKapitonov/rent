package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	pathToConfig = "/.config/rent/setting.yaml"
)

var (
	ErrPathNotExist    = errors.New("path does not exist")
	ErrPathNotIsHome   = errors.New("path is not in the home directory")
	ErrPathNotAbsolute = errors.New("path is not absolute")
	ErrPathNotDir      = errors.New("path is not a directory")
	ErrOneArgument     = errors.New("maximum one argument")
)

var settingCmd = &cobra.Command{
	Use:   "setting",
	Short: "Set default output directory for torrent files.",
	Long: `Set default output directory for torrent files.
Write command: rent setting </absolute/path/to/output_directory/>.
By default it is $HOME directory`,
	Args:    cobra.MaximumNArgs(1),
	Run:     SetOutDir,
	Example: "rent setting /home/Video/",
}

// SetOutDir sets the output directory for moving files.
func SetOutDir(_ *cobra.Command, args []string) {
	setConfigFile()

	isOne, err := isOneArgument(args)
	cobra.CheckErr(err)
	if !isOne {
		cobra.WriteStringAndCheck(
			stderr,
			fmt.Sprintf(
				"Current path: %s \nFor change write:\n\trent settings </home/$USER/path>",
				viper.GetString("out_dir")),
		)
		return
	}

	outPath := args[0]

	outPath = filepath.Clean(outPath)

	if err := validateOutputPath(outPath); err != nil {
		cobra.CheckErr(err)
	}

	outPath += "/"

	viper.Set("out_dir", outPath)

	if err := viper.WriteConfig(); err != nil {
		cobra.CheckErr(err)
	}

	cobra.WriteStringAndCheck(
		stderr,
		"Output directory has been successfully selected :)",
	)
}

func isOneArgument(args []string) (bool, error) {
	if len(args) == 1 && args[0] != "" {
		return true, nil
	}

	if len(args) > 1 {
		return false, ErrOneArgument
	}

	return false, nil
}

// validateOutputPath checks if the provided path exists and is a directory.
func validateOutputPath(path string) error {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return ErrPathNotExist
	}

	if !strings.HasPrefix(path, viper.GetString("HOME")) {
		return ErrPathNotIsHome
	}

	if !stat.IsDir() {
		return ErrPathNotDir
	}

	return nil
}

// setConfigFile sets the config file name for writing the output directory.
func setConfigFile() {
	viper.AutomaticEnv()
	homeEnv := viper.GetString("HOME")
	viper.SetConfigFile(
		homeEnv + pathToConfig,
	)

	if err := viper.ReadInConfig(); err != nil {
		cobra.CheckErr(err)
	}
}
