package constant

import (
	"log"
	"os"
	"path/filepath"
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
		if dir, got := getFirstValidEnv("NALI_CONFIG_HOME", "XDG_CONFIG_HOME"); got {
			ConfigDirPath = dir
		} else {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				log.Fatal("get user's home dir failed!")
			}
			ConfigDirPath = filepath.Join(homeDir, ".nali")
		}

		if dir, got := getFirstValidEnv("NALI_DB_HOME", "XDG_DATA_HOME"); got {
			DataDirPath = dir
		} else {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				log.Fatal("get user's home dir failed!")
			}
			DataDirPath = filepath.Join(homeDir, ".nali")
		}
	}

	prepareDir(ConfigDirPath)
	prepareDir(DataDirPath)

	os.Chdir(DataDirPath)
}

func getFirstValidEnv(keys ...string) (string, bool) {
	for _, key := range keys {
		if value := os.Getenv(key); len(value) > 0 {
			return value, true
		}
	}
	return "", false
}

func prepareDir(dir string) {
	stat, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatal("can not create config dir:", dir)
		}
	}
	if !stat.IsDir() {
		log.Fatal("path already exists, but not a dir:", dir)
	}
}
