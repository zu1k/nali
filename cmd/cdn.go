package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/zu1k/nali/internal/app"

	"github.com/spf13/cobra"
)

// cdnCmd represents the cdn command
var cdnCmd = &cobra.Command{
	Use:   "cdn",
	Short: "Query cdn service provider",
	Long:  `Query cdn service provider`,
	Run: func(cmd *cobra.Command, args []string) {
		if update {
			app.UpdateDB()
		}

		app.InitCDNDB()

		if len(args) == 0 {
			stdin := bufio.NewScanner(os.Stdin)
			for stdin.Scan() {
				line := stdin.Text()
				if line == "quit" || line == "exit" {
					return
				}
				fmt.Println(app.ReplaceCDNInString(line))
			}
		} else {
			app.ParseCDNs(args)
		}
	},
}

var (
	update = false
)

func init() {
	rootCmd.AddCommand(cdnCmd)
	cdnCmd.Flags().BoolVarP(&update, "update", "u", false, "Update CDN database")
}
