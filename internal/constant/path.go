package constant

import (
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

var (
	ConfigDirPath string
	DataDirPath   string
)

func init() {
	if naliHome := os.Getenv("NALI_HOME"); len(naliHome) != 0 {
		ConfigDirPath = naliHome
		DataDirPath = naliHome
	} else {
		ConfigDirPath = os.Getenv("NALI_CONFIG_HOME")
		if len(ConfigDirPath) == 0 {
			ConfigDirPath = filepath.Join(xdg.ConfigHome, "nali")
		}

		DataDirPath = os.Getenv("NALI_DB_HOME")
		if len(DataDirPath) == 0 {
			DataDirPath = filepath.Join(xdg.DataHome, "nali")
		}
	}

	prepareDir(ConfigDirPath)
	prepareDir(DataDirPath)

	_ = os.Chdir(DataDirPath)
}

func prepareDir(dir string) {
	stat, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatal("can not create config dir:", dir)
		}
	} else if err != nil {
		log.Fatal(err)
	} else if !stat.IsDir() {
		log.Fatal("path already exists, but not a dir:", dir)
	}
}
