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

var stderr = os.Stderr

var downlaodCmd = &cobra.Command{
	Use:   "download",
	Short: `Designed for installing and downloading a torrent file.`,
	Long: `Designed for installing and downloading a torrent file. When using this command,
	Rent will download the specified torrent file to your computer.`,
	Run:               download,
	Args:              cobra.MinimumNArgs(1),
	Example:           "rent download path/to/<some_name>.torrent",
	ValidArgsFunction: completeTorrentFiles,
}

// completeTorrentFiles returns files ending on ".torrent" in current directory.
func completeTorrentFiles(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	files, err := filepath.Glob("*.torrent")
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	return files, cobra.ShellCompDirectiveDefault
}

// download cobra command for downloading torrent files.
func download(cmd *cobra.Command, args []string) {
	pathToFile, err := getTfilePath(args)
	cobra.CheckErr(err)

	client, err := torrent.NewClient(nil)
	cobra.CheckErr(err)

	tfile, err := client.AddTorrentFromFile(pathToFile)
	cobra.CheckErr(err)

	<-tfile.GotInfo()

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

	go loading(tfile, bar)
	tfile.DownloadAll()

	go shutdown(client)

	done := client.WaitAll()

	fileName := tfile.Info().Name
	err = move(pathToFile, fileName, done)

	cobra.CheckErr(err)
}

func getTfilePath(args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("Invalid argument")
	}
	filePath := args[0]

	arr := strings.Split(filePath, ".")
	if arr[len(arr)-1] != "torrent" {
		return "", errors.New("It is not torrent file!")
	}

	return filePath, nil
}

func shutdown(client *torrent.Client) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch

	client.Closed()
	client.Close()
}

func move(pathToFile string, fileName string, done bool) error {
	if done {
		fmt.Fprintln(stderr, "\nFile successfily downloaded :)")
	} else {
		fmt.Fprintln(stderr, "\n Download was paused ")
	}
	newPath := getNewPath(fileName)
	currentPath := getCurrentPath(fileName)

	if newPath == "" {
		return nil
	}

	return rename(currentPath, newPath)
}

// getCurrentPath returns absolute path of files from torrent file.
func getCurrentPath(name string) string {

	currentPath, err := filepath.Abs(filepath.Base(name))
	cobra.CheckErr(err)
	return currentPath
}

// getNewPath returns path from .config/settings.yaml file + name of file.
func getNewPath(name string) string {
	setConfigFile()

	return viper.GetString("out_dir") + name
}

// rename move file with check existing .
func rename(oldPath string, newPath string) error {
	err := os.Rename(oldPath, newPath)
	if err != nil {
		newPath = checkExist(err, newPath)
		err = os.Rename(oldPath, newPath)
	}

	return err
}

// checkExist is check if file or dir exist and add Unix time for unicue name.
func checkExist(err error, newPath string) string {
	defPath := viper.GetString("out_dir")

	if errors.Is(err, os.ErrExist) {
		newPath = getUniqueName(newPath, defPath)
	}

	return newPath
}

// getUniqueName returns unique file name for dircetory.
func getUniqueName(newPath string, defPath string) string {
	var indx int
	arr := strings.Split(filepath.Base(newPath), ".")
	if len(arr) >= 2 {
		indx = len(arr) - 2
	} else {
		indx = 0
	}

	arr[indx] += fmt.Sprintf("_(%s)", time.Now().Format(time.DateTime))

	newPath = defPath + "/" + strings.Join(arr, ".")
	return newPath
}

// loading add points to progress bar every 50 Millisecond.
func loading(tfile *torrent.Torrent, bar *progressbar.ProgressBar) {
	fmt.Fprintf(stderr, "\n Start downloading: %s \n", tfile.Name())
	for !tfile.Complete.Bool() {
		current := tfile.Stats().BytesReadData
		time.Sleep(50 * time.Millisecond)
		afterSecond := tfile.Stats().BytesReadData
		if err := bar.Add64(afterSecond.Int64() - current.Int64()); err != nil {
			if errors.Is(err, errors.New("current number exceeds max")) {
				bar.ChangeMax((int(float64(bar.GetMax()) * 1.1))) // Increases max length of bar by 10%.
			}
			cobra.CheckErr(err)
		}
	}
}
