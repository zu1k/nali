package cmd

import (
	"github.com/zu1k/nali/internal/db"
	"strings"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "update chunzhen, zxipv6, ip2region ip database and cdn",
	Long:    `update chunzhen, zxipv6, ip2region ip database and cdn. Use commas to separate`,
	Example: "nali update --db chunzhen,cdn",
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
	rootCmd.AddCommand(updateCmd)
	rootCmd.PersistentFlags().String("db", "", "choose db you want to update")
}
