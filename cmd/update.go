package cmd

import (
	"log"
	"strings"

	"github.com/zu1k/nali/internal/db"
	"github.com/zu1k/nali/internal/repo"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:     "update [--db dbs -v]",
	Short:   "update qqwry, zxipv6wry, ip2region ip database and cdn, update nali to latest version if -v",
	Long:    `update qqwry, zxipv6wry, ip2region ip database and cdn. Use commas to separate. update nali to latest version if -v`,
	Example: "nali update --db qqwry,cdn -v",
	Run: func(cmd *cobra.Command, args []string) {
		DBs, _ := cmd.Flags().GetString("db")

		version, _ := cmd.Flags().GetBool("v")
		if version {
			if err := repo.UpdateRepo(); err != nil {
				log.Printf("update nali to latest version failed: %v \n", err)
			}
		}

		var DBNameArray []string
		if DBs != "" {
			DBNameArray = strings.Split(DBs, ",")
		}
		db.UpdateDB(DBNameArray...)
	},
}

func init() {
	updateCmd.PersistentFlags().String("db", "", "choose db you want to update")
	updateCmd.PersistentFlags().Bool("v", false, "decide whether to update the nali version")
	rootCmd.AddCommand(updateCmd)
}
