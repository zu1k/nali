package cmd

import (
	"log"
	"path/filepath"

	"github.com/zu1k/nali/internal/constant"

	"github.com/zu1k/nali/pkg/cdn"
	"github.com/zu1k/nali/pkg/zxipv6wry"

	"github.com/zu1k/nali/pkg/qqwry"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update chunzhen ip database",
	Long:  `update chunzhen ip database`,
	Run: func(cmd *cobra.Command, args []string) {
		// Chunzhen ipv4
		filePath := filepath.Join(constant.HomePath, "qqwry.dat")
		log.Println("正在下载最新 纯真 IPv4数据库...")
		_, err := qqwry.Download(filePath)
		if err != nil {
			log.Fatalln("下载失败", err.Error())
			return
		}

		// ZX ipv6
		filePath = filepath.Join(constant.HomePath, "ipv6wry.database")
		log.Println("正在下载最新 ZX IPv6数据库...")
		_, err = zxipv6wry.Download(filePath)
		if err != nil {
			log.Fatalln("下载失败", err.Error())
			return
		}

		// cdn
		filePath = filepath.Join(constant.HomePath, "cdn.json")
		log.Println("正在下载最新 CDN服务提供商数据库...")
		_, err = cdn.Download(filePath)
		if err != nil {
			log.Fatalln("下载失败", err.Error())
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
