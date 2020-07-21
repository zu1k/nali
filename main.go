package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/zu1k/nali/cmd"
	"github.com/zu1k/nali/constant"
)

func main() {
	setHomePath()
	cmd.Execute()
}

func setHomePath() {
	homePath := os.Getenv("NALI_DB_HOME")
	if homePath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		homePath = filepath.Join(homeDir, ".nali")
	}
	constant.HomePath = homePath
	if _, err := os.Stat(homePath); os.IsNotExist(err) {
		if err := os.MkdirAll(homePath, 0777); err != nil {
			log.Fatal("can not create", homePath, ", use bin dir instead")
		}
	}
}
