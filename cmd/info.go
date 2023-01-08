package cmd

import (
	"fmt"

	"github.com/zu1k/nali/internal/constant"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:     "info",
	Short:   "get the necessary information of nali",
	Long:    `get the necessary information of nali`,
	Example: "nali info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Nali Version:     ", constant.Version)
		fmt.Println("Config Dir Path:  ", constant.ConfigDirPath)
		fmt.Println("DB Data Dir Path: ", constant.DataDirPath)

		fmt.Println("Selected IPv4 DB: ", viper.GetString("selected.ipv4"))
		fmt.Println("Selected IPv6 DB: ", viper.GetString("selected.ipv6"))
		fmt.Println("Selected CDN DB:  ", viper.GetString("selected.cdn"))
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
