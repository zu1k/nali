package constant

import (
	"log"
	"os"
	"path/filepath"
)

var (
	// WorkDirPath database home path
	WorkDirPath string
)

func init() {
	WorkDirPath = os.Getenv("NALI_HOME")
	if WorkDirPath == "" {
		WorkDirPath = os.Getenv("NALI_DB_HOME")
	}
	if WorkDirPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		WorkDirPath = filepath.Join(homeDir, ".nali")
	}
	if _, err := os.Stat(WorkDirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(WorkDirPath, 0777); err != nil {
			log.Fatal("can not create", WorkDirPath, ", use bin dir instead")
		}
	}

	_ = os.Chdir(WorkDirPath)
}
