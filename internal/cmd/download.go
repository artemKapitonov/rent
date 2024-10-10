package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const timeOfUpdatingSpeed = 50 * time.Millisecond
const increaseFiveProcent = 1.05

var stderr = os.Stderr

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: `Designed for installing and downloading a torrent file`,
	Long: `Designed for installing and downloading a torrent file. When using this command,
rent will download the specified torrent file to your computer`,
	Run:               download,
	Args:              onlyOneArg,
	Example:           "rent download path/to/<some_name>.torrent",
	ValidArgsFunction: completeTorrentFiles,
}

func onlyOneArg(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("command accepts only 1 argument, no less and no more")
	}

	return nil
}

// completeTorrentFiles returns files ending on ".torrent" in current directory.
func completeTorrentFiles(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	files, err := filepath.Glob("*.torrent")
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	return files, cobra.ShellCompDirectiveDefault
}

// download cobra command for downloading torrent files.
func download(_ *cobra.Command, args []string) {
	pathToFile, err := getTfilePath(args)
	cobra.CheckErr(err)

	client, err := torrent.NewClient(nil)
	cobra.CheckErr(err)

	tfile, err := client.AddTorrentFromFile(pathToFile)
	cobra.CheckErr(err)

	<-tfile.GotInfo()

	bar := createProgressBar(tfile)

	go loading(tfile, bar)
	tfile.DownloadAll()

	go shutdown(client)

	done := client.WaitAll()
	if done {
		cobra.WriteStringAndCheck(stderr, "Download was completed successfully :)")
	}

	fileName := tfile.Info().Name
	err = move(fileName, done)
	cobra.CheckErr(err)
}

func createProgressBar(tfile *torrent.Torrent) *progressbar.ProgressBar {
	bar := progressbar.NewOptions(
		int(tfile.Info().TotalLength()),
		progressbar.OptionSetWriter(stderr),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]·[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	return bar
}

func getTfilePath(args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("invalid argument: expected exactly one argument")
	}

	filePath := args[0]
	if !isTorrentFile(filePath) {
		return "", errors.New("invalid file type: expected a .torrent file")
	}

	return filePath, nil
}

// isTorrentFile checks if the given file path has a .torrent extension.
func isTorrentFile(filePath string) bool {
	return strings.HasSuffix(filePath, ".torrent")
}

func shutdown(client *torrent.Client) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	<-ch
	client.Closed()

	if err := errors.Join(client.Close()...); err != nil {
		cobra.WriteStringAndCheck(stderr, err.Error())
	}
}

func move(fileName string, done bool) error {
	if done {
		cobra.WriteStringAndCheck(stderr, "\nFile successfully downloaded :)")
	} else {
		cobra.WriteStringAndCheck(stderr, "\nDownload was paused.")
	}

	newPath := getNewPath(fileName)
	currentPath := getCurrentPath(fileName)

	if newPath == "" {
		return nil
	}

	return rename(currentPath, newPath)
}

// getCurrentPath returns the absolute path of the downloaded file.
func getCurrentPath(name string) string {
	currentPath, err := filepath.Abs(filepath.Base(name))
	cobra.CheckErr(err)

	return currentPath
}

// getNewPath constructs the new path from the config and file name.
func getNewPath(name string) string {
	setConfigFile()
	return filepath.Join(viper.GetString("out_dir"), name)
}

// rename moves the file to the new path, checking for existing files.
func rename(oldPath, newPath string) error {
	newPath = getUniqueName(newPath)

	err := os.Rename(oldPath, newPath)
	if err != nil {
		return errors.New("can't move file to output directory")
	}

	return nil
}

// getUniqueName generates a unique file name for the directory.
func getUniqueName(newPath string) string {
	defPath := viper.GetString("out_dir")
	baseName := filepath.Base(newPath)

	arr := strings.Split(baseName, ".")

	// Determine the index to append the timestamp
	indx := len(arr) - 2
	if len(arr) < 2 {
		indx = 0
	}

	arr[indx] += fmt.Sprintf(
		"_(%s)",
		time.Now().Format("2006-01-02_15-04-05"),
	)

	return filepath.Join(defPath, strings.Join(arr, "."))
}

// loading updates the progress bar every 50 milliseconds.
func loading(tfile *torrent.Torrent, bar *progressbar.ProgressBar) {
	cobra.WriteStringAndCheck(
		stderr,
		"Start downloading: %s"+tfile.Name(),
	)

	for !tfile.Complete.Bool() {
		current := tfile.Stats().BytesReadData

		time.Sleep(timeOfUpdatingSpeed)

		afterTime := tfile.Stats().BytesReadData
		bytesRead := afterTime.Int64() - current.Int64()

		if err := bar.Add64(bytesRead); err != nil {
			if errors.Is(err, errors.New("current number exceeds max")) {
				bar.ChangeMax(int(float64(bar.GetMax()) * increaseFiveProcent)) // Increase max length of bar by 5%.
				continue
			}

			cobra.WriteStringAndCheck(stderr, err.Error())
		}
	}
}
