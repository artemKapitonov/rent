package cmd

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var stderr = os.Stderr

var validArgs = []string{"download", "settings"}

var downlaodCmd = &cobra.Command{
	Use: "download",
	Short: `Предназначена для установки и загрузки торрент-файла.
При использовании этой команды Rent будет загружать указанный торрент-файл на ваш компьютер.`,
	Long: `Предназначена для установки и загрузки торрент-файла с именем, указанным в параметре -f.
	При использовании этой команды Rent будет загружать указанный торрент-файл на ваш компьютер.`,
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
		fmt.Fprintln(stderr, "File successfily downloaded :)")
	} else {
		fmt.Fprintln(stderr, "\n Download was paused ")
	}
	newPath := getNewPath(pathToFile, fileName)
	currentPath := getCurrentPath(fileName)

	return rename(currentPath, newPath)
}

// getCurrentPath returns absolute path of files from torrent file.
func getCurrentPath(name string) string {

	currentPath, err := filepath.Abs(filepath.Base(name))
	cobra.CheckErr(err)
	return currentPath
}

// TODO: add defaul path from config
func getNewPath(pathToFile string, name string) string {
	arr := strings.Split(pathToFile, "/")
	arr[len(arr)-1] = name

	newPath := "/home/kapitonov/Видео/rent/" + name

	return newPath
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

// TODO: remake with path from .config/rent.yaml
// checkExist is check if file or dir exist and add (1) for unicue name.
func checkExist(err error, newPath string) string {
	if errors.Is(err, os.ErrExist) {
		arr := strings.Split(filepath.Base(newPath), ".")
		arr[len(arr)-2] += "." + strconv.Itoa(rand.Intn(1000))
		newPath = "/home/kapitonov/Видео/rent/" + strings.Join(arr, ".")
	}

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
			cobra.CheckErr(err)
		}
	}
}
