package cmd

import (
	"github.com/zu1k/nali/internal/app"
	"github.com/zu1k/nali/internal/ipdb"

	"github.com/spf13/cobra"
)

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Query IP information",
	Long: `Query IP information.

#1 Query a simple IP address

	$ nali 1.2.3.4

  or use pipe

	$ echo IP 6.6.6.6 | nali

#2 Query multiple IP addresses

	$ nali 1.2.3.4 4.3.2.1 123.23.3.0

#3 Interactive query

	$ nali
	123.23.23.23
	123.23.23.23 [越南 越南邮电集团公司]
	quit

#4 Use with dig

	$ dig nali.lgf.im +short | nali

#5 Use with nslookup

	$ nslookup nali.lgf.im 8.8.8.8 | nali

#6 Use with any other program

	bash abc.sh | nali

#7 IPV6 support

	$ nslookup google.com | nali
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitIPDB(ipdb.GetIPDBType())
		app.ParseIPs(args)
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)
}
