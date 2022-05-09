package main

import (
	"github.com/zu1k/nali/cmd"
	"github.com/zu1k/nali/internal/config"
	"github.com/zu1k/nali/internal/constant"
)

func main() {
	config.ReadConfig(constant.WorkDirPath)
	cmd.Execute()
}
