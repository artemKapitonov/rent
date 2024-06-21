package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/anacrolix/torrent"
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
	Short: `Предназначена для установки и загрузки торрент-файла с именем, указанным в параметре -f.
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

	go loading(tfile, client)

	tfile.DownloadAll()

	client.WaitAll()

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

func loading(tfile *torrent.Torrent, client *torrent.Client) {
	for {
		time.Sleep(2 * time.Second)
		total := tfile.Info().TotalLength()
		current := tfile.Stats().BytesReadData
		fmt.Println(total, "total")
		fmt.Println(current, "current")

		proc := (float64(current.Int64()) / float64(total)) * 100
		time.Sleep(1 * time.Second)
		fmt.Println(proc)
	}

}
