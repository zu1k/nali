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
	Use:     "update [--db dbs] [--version v]",
	Short:   "update qqwry, zxipv6wry, ip2region ip database and cdn",
	Long:    `update qqwry, zxipv6wry, ip2region ip database and cdn. Use commas to separate`,
	Example: "nali update --db qqwry,cdn",
	Run: func(cmd *cobra.Command, args []string) {
		DBs, _ := cmd.Flags().GetString("db")
		version, _ := cmd.Flags().GetString("version")

		if err := repo.UpdateRepo(version); err != nil {
			log.Printf("update nali to version %s failed: %v", version, err)
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
	updateCmd.PersistentFlags().String("version", "", "choose nali version you want to update")
	rootCmd.AddCommand(updateCmd)
}
