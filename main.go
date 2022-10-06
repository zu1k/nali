package main

import (
	"github.com/zu1k/nali/internal/constant"

	"github.com/zu1k/nali/cmd"
	"github.com/zu1k/nali/internal/config"

	_ "github.com/zu1k/nali/internal/migration"
)

func main() {
	config.ReadConfig(constant.ConfigDirPath)
	cmd.Execute()
}
