package app

import (
	"bufio"
	"os"
	"runtime"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var needTransform = false

func init() {
	stat, _ := os.Stdin.Stat()
	needTransform = ((stat.Mode() & os.ModeNamedPipe) != 0) && runtime.GOOS == "windows"
}

func Root(args []string) {
	if len(args) == 0 {
		stdin := bufio.NewScanner(os.Stdin)
		for stdin.Scan() {
			line := stdin.Text()
			if needTransform {
				line, _, _ = transform.String(simplifiedchinese.GBK.NewDecoder(), line)
			}
			if line == "quit" || line == "exit" {
				return
			}
			// TODO: pring line
			//fmt.Printf("%s\n", ReplaceIPInString(ReplaceCDNInString(line)))
		}
	} else {
		// TODO: do something
		//ParseIPs(args)
	}
}
