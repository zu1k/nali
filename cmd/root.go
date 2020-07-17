package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/zu1k/nali/internal/app"
	"github.com/zu1k/nali/internal/ipdb"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nali",
	Short: "An offline tool for querying IP geographic information",
	Long:  `An offline tool for querying IP geographic information.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitIPDB(ipdb.GetIPDBType())
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
