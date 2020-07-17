package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/zu1k/nali/pkg/qqwry"

	"github.com/zu1k/nali/constant"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update chunzhen ip database",
	Long:  `update chunzhen ip database`,
	Run: func(cmd *cobra.Command, args []string) {
		filePath := filepath.Join(constant.HomePath, "qqwry.dat")

		log.Println("正在下载最新纯真 IP 库...")
		tmpData, err := qqwry.GetOnline()
		if err != nil {
			log.Fatalln("下载失败", err.Error())
			return
		}

		// 文件存在就删除
		_, err = os.Stat(filePath)
		if err == nil {
			err = os.Remove(filePath)
			if err != nil {
				log.Fatalln("旧文件删除失败", err.Error())
				os.Exit(1)
			}
		}

		if err := ioutil.WriteFile(filePath, tmpData, 0644); err == nil {
			log.Printf("已将最新的纯真 IP 库保存到本地 %s ", filePath)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
