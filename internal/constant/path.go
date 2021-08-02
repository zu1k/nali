package constant

import (
	"log"
	"os"
	"path/filepath"
)

var (
	// HomePath database home path
	HomePath string
)

func init() {
	HomePath = os.Getenv("NALI_DB_HOME")
	if HomePath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		HomePath = filepath.Join(homeDir, ".nali")
	}
	if _, err := os.Stat(HomePath); os.IsNotExist(err) {
		if err := os.MkdirAll(HomePath, 0777); err != nil {
			log.Fatal("can not create", HomePath, ", use bin dir instead")
		}
	}
}
