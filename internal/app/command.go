package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/zu1k/nali/internal/entity"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func Root(args []string, needTransform bool) {
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
			fmt.Printf("%s\n", entity.ParseLine(line).ColorString())
		}
	} else {
		fmt.Printf("%s\n", entity.ParseLine(strings.Join(args, " ")).ColorString())
	}
}
