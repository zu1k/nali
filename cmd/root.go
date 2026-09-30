package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

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
	RunE: func(cmd *cobra.Command, args []string) error {
		gbk, _ := cmd.Flags().GetBool("gbk")
		isJson, _ := cmd.Flags().GetBool("json")
		colorMode, _ := cmd.Flags().GetString("color")

		switch colorMode {
		case "always":
			color.NoColor = false
		case "never":
			color.NoColor = true
		case "auto":
			// fatih/color already disables colors when stdout is not a
			// terminal, NO_COLOR is set or TERM=dumb
		default:
			return fmt.Errorf("invalid --color value %q: must be auto, always or never", colorMode)
		}

		if len(args) == 0 {
			autoGBK := common.ConsoleUsesGBK()
			stdin := bufio.NewScanner(os.Stdin)
			stdin.Split(common.ScanLines)
			for stdin.Scan() {
				line := common.DecodeInput(stdin.Text(), gbk, autoGBK)
				if line := strings.TrimSpace(line); line == "quit" || line == "exit" {
					return nil
				}
				if isJson {
					_, _ = fmt.Fprintf(color.Output, "%s", entity.ParseLine(line).Json())
				} else {
					_, _ = fmt.Fprintf(color.Output, "%s", entity.ParseLine(line).ColorString())
				}
			}
		} else {
			if isJson {
				_, _ = fmt.Fprintf(color.Output, "%s", entity.ParseLine(strings.Join(args, " ")).Json())
			} else {
				for _, line := range args {
					_, _ = fmt.Fprintf(color.Output, "%s\n", entity.ParseLine(line).ColorString())
				}
			}
		}
		return nil
	},
}

// Execute parse subcommand and run
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}

func init() {
	rootCmd.Flags().Bool("gbk", false, "Decode input as GBK (detected automatically on Chinese Windows)")
	rootCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
	rootCmd.Flags().String("color", "auto", "Colorize output: auto, always or never (NO_COLOR is honored in auto mode)")
}
