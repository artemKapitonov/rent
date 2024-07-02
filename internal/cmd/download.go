package cmd

import (
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
)

var downlaodCmd = &cobra.Command{
	Use: "download",
	Short: `Предназначена для установки и загрузки торрент-файла.
При использовании этой команды Rent будет загружать указанный торрент-файл на ваш компьютер.`,
	Long: `Предназначена для установки и загрузки торрент-файла с именем, указанным в параметре -f.
	При использовании этой команды Rent будет загружать указанный торрент-файл на ваш компьютер.`,
	Run:     Download,
	Example: "rent download <some_name>.torrent",
}

func Download(cmd *cobra.Command, args []string) {
	pathToFile := args[0]
	err := recover()
	cobra.CheckErr(err)

	fmt.Println(pathToFile)

	client, err := torrent.NewClient(nil)
	cobra.CheckErr(err)

	tfile, err := client.AddTorrentFromFile(pathToFile)
	cobra.CheckErr(err)

	fmt.Println("file founded")
	<-tfile.GotInfo()

	bar := progressbar.NewOptions(
		int(tfile.Info().TotalLength()),
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionShowBytes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "·",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	done := make(chan bool, 1)
	go loading(tfile, client, bar)
	tfile.DownloadAll()
	go WaitLoading(client, done)

	Move(pathToFile, tfile, done)
}

func WaitLoading(client *torrent.Client, done chan bool) {
	fmt.Println("wait download")
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-ch:
		done <- false
	}

	done <- client.WaitAll()
}

func Move(pathToFile string, tfile *torrent.Torrent, done chan bool) {
	if <-done {
		fmt.Println("File successfily downloaded")
	} else {
		fmt.Println("Download was paused")
	}

	arr := strings.Split(pathToFile, "/")
	arr[len(arr)-1] = tfile.Info().Name

	newPath, err := filepath.Abs("/home/kapitonov/Видео/films/" + tfile.Info().Name)
	cobra.CheckErr(err)

	oldPath, err := filepath.Abs(filepath.Base(strings.Join(arr, "/")))
	cobra.CheckErr(err)

	fmt.Println(oldPath, "\n", newPath)

	err = os.Rename(oldPath, newPath)

	cobra.CheckErr(err)
	fmt.Println("file move success")
}

func loading(tfile *torrent.Torrent, client *torrent.Client, bar *progressbar.ProgressBar) {
	for {
		current := tfile.Stats().BytesReadData
		time.Sleep(3 * time.Millisecond)
		afterSecond := tfile.Stats().BytesReadData
		bar.Add(int(afterSecond.Int64() - current.Int64()))
	}
}
