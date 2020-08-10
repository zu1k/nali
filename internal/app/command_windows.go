// +build windows

package app

import (
	"bufio"
	"fmt"
	"os"
	"syscall"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var pipeline = false

func init() {
	ftype, _ := syscall.GetFileType(syscall.Handle(os.Stdin.Fd()))
	pipeline = ftype == 3
}

func Root(args []string) {
	if len(args) == 0 {
		stdin := bufio.NewScanner(os.Stdin)
		for stdin.Scan() {
			line := stdin.Text()
			if pipeline {
				line, _, _ = transform.String(simplifiedchinese.GBK.NewDecoder(), line)
			}
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
			if pipeline {
				line, _, _ = transform.String(simplifiedchinese.GBK.NewDecoder(), line)
			}
			if line == "quit" || line == "exit" {
				return
			}
			fmt.Println(ReplaceCDNInString(line))
		}
	} else {
		ParseCDNs(args)
	}
}
