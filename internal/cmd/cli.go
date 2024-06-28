package cmd

import (
	"os"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

type service interface {
	Download(path string) error
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rent",
	Short: "Rent - это легковесный CLI торрент-клиент.",
	Long: `Rent - это легковесный CLI торрент-клиент для быстрого и удобного скачивания файлов напрямую из терминала.
С помощью Rent вы можете легко управлять загрузками,
просматривать статус загрузок и контролировать скорость скачивания, не покидая командную строку.

Наслаждайтесь быстрым и эффективным скачиванием файлов с помощью Rent!"`,
}

var downlaodCmd = &cobra.Command{
	Use: "download",
	Short: `Предназначена для установки и загрузки торрент-файла.
При использовании этой команды Rent будет загружать указанный торрент-файл на ваш компьютер.`,
	Long: `Предназначена для установки и загрузки торрент-файла с именем, указанным в параметре -f.
	При использовании этой команды Rent будет загружать указанный торрент-файл на ваш компьютер.`,
	Run: Download,
}

func Download(cmd *cobra.Command, args []string) {
	pathToFile := os.Args[2]

	client, err := torrent.NewClient(nil)
	cobra.CheckErr(err)

	tfile, err := client.AddTorrentFromFile(pathToFile)
	cobra.CheckErr(err)

	<-tfile.GotInfo()

	bar := progressbar.NewOptions(
		int(tfile.Info().TotalLength()),
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionShowBytes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetDescription(tfile.Info().BestName()),
		progressbar.OptionFullWidth(),
	)

	go loading(tfile, client, bar)

	tfile.DownloadAll()

	client.WaitAll()


	arr := strings.Split(pathToFile, "/")
	newPath := "/home/kapitonov/Видео/films/" + tfile.Info().BestName()
	arr[len(arr)-1] = tfile.Info().BestName()
	err = os.Rename(strings.Join(arr, "/"), newPath)
	cobra.CheckErr(err)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(downlaodCmd)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func loading(tfile *torrent.Torrent, client *torrent.Client, bar *progressbar.ProgressBar) {

	for {
		current := tfile.Stats().BytesReadData
		time.Sleep(20 * time.Millisecond)
		afterSecond := tfile.Stats().BytesReadData
		bar.Add(int(afterSecond.Int64() - current.Int64()))
	}
}
