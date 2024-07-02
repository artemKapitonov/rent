package cmd

import (
	"os"

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

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(0)
	}
}

func init() {
	rootCmd.AddCommand(downlaodCmd)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
