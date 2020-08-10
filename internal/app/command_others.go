// +build !windows

package app

import (
	"bufio"
	"fmt"
	"os"
)

func Root(args []string) {
	if len(args) == 0 {
		stdin := bufio.NewScanner(os.Stdin)
		for stdin.Scan() {
			line := stdin.Text()
			if line == "quit" || line == "exit" {
				return
			}
			fmt.Printf("%s\n", ReplaceIPInString(ReplaceCDNInString(line)))
		}
	} else {
		ParseIPs(args)
	}
}

func CDN(args []string) {
	if len(args) == 0 {
		stdin := bufio.NewScanner(os.Stdin)
		for stdin.Scan() {
			line := stdin.Text()
			if line == "quit" || line == "exit" {
				return
			}
			fmt.Println(ReplaceCDNInString(line))
		}
	} else {
		ParseCDNs(args)
	}
}
