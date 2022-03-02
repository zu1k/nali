package cmd

import (
	"github.com/zu1k/nali/internal/db"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update chunzhen ip database",
	Long:  `update chunzhen ip database`,
	Run: func(cmd *cobra.Command, args []string) {
		db.UpdateAllDB()

	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
