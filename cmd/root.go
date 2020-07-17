package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/zu1k/nali/internal/app"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nali",
	Short: "",
	Long:  ``,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			stdin := bufio.NewScanner(os.Stdin)
			for stdin.Scan() {
				line := stdin.Text()
				if line == "quit" || line == "exit" {
					return
				}
				fmt.Println(app.ReplaceInString(line))
			}
		} else {
			app.ParseIPs(args)
		}
	},
}

// Execute parse subcommand and run
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
