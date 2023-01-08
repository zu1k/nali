package cmd

import (
	"strings"

	"github.com/zu1k/nali/internal/db"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "update qqwry, zxipv6wry, ip2region ip database and cdn",
	Long:    `update qqwry, zxipv6wry, ip2region ip database and cdn. Use commas to separate`,
	Example: "nali update --db qqwry,cdn",
	Run: func(cmd *cobra.Command, args []string) {
		DBs, _ := cmd.Flags().GetString("db")
		var DBNameArray []string
		if DBs != "" {
			DBNameArray = strings.Split(DBs, ",")
		}
		db.UpdateDB(DBNameArray...)
	},
}

func init() {
	updateCmd.PersistentFlags().String("db", "", "choose db you want to update")
	rootCmd.AddCommand(updateCmd)
}
