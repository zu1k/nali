package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"github.com/zu1k/nali/internal/constant"
	"github.com/zu1k/nali/pkg/common"
	"github.com/zu1k/nali/pkg/entity"
)

var rootCmd = &cobra.Command{
	Use:   "nali",
	Short: "An offline tool for querying IP geographic information",
	Long: `An offline tool for querying IP geographic information.

Find document on: https://github.com/zu1k/nali

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

	$ dig nali.zu1k.com +short | nali

#5 Use with nslookup

	$ nslookup nali.zu1k.com 8.8.8.8 | nali

#6 Use with any other program

	bash abc.sh | nali

#7 IPV6 support
`,
	Version: constant.Version,
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		gbk, _ := cmd.Flags().GetBool("gbk")

		if len(args) == 0 {
			stdin := bufio.NewScanner(os.Stdin)
			stdin.Split(common.ScanLines)
			for stdin.Scan() {
				line := stdin.Text()
				if gbk {
					line, _, _ = transform.String(simplifiedchinese.GBK.NewDecoder(), line)
				}
				if line := strings.TrimSpace(line); line == "quit" || line == "exit" {
					return
				}
				_, _ = fmt.Fprintf(color.Output, "%s", entity.ParseLine(line).ColorString())
			}
		} else {
			_, _ = fmt.Fprintf(color.Output, "%s\n", entity.ParseLine(strings.Join(args, " ")).ColorString())
		}
	},
}

// Execute parse subcommand and run
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}

func init() {
	rootCmd.Flags().Bool("gbk", false, "Use GBK decoder")
}
