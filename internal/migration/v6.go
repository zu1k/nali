package migration

import (
	"os"
	"path/filepath"

	"github.com/google/martian/log"
	"github.com/zu1k/nali/internal/constant"
)

func migration2v6() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	oldDefaultWorkPath := filepath.Join(homeDir, ".nali")

	oldDefaultWorkPath, err = filepath.Abs(oldDefaultWorkPath)
	if err != nil {
		log.Errorf("Get absolute path for oldDefaultWorkPath failed: %s\n", err)
	}
	mewWorkPath, err := filepath.Abs(constant.ConfigDirPath)
	if err != nil {
		log.Errorf("Get absolute path for mewWorkPath failed: %s\n", err)
	}
	if oldDefaultWorkPath == mewWorkPath {
		// User chooses to continue using old directory
		return
	}

	_, err = os.Stat(oldDefaultWorkPath)
	if err == nil {
		println("Old data directories are detected and will attempt to migrate automatically")

		oldDefaultConfigPath := filepath.Join(oldDefaultWorkPath, "config.yaml")
		stat, err := os.Stat(oldDefaultConfigPath)
		if err == nil {
			if stat.Mode().IsRegular() {
				_ = os.Rename(oldDefaultConfigPath, filepath.Join(constant.ConfigDirPath, "config.yaml"))
			}
		}

		files, err := os.ReadDir(oldDefaultWorkPath)
		if err == nil {
			for _, file := range files {
				if file.Type().IsRegular() {
					_ = os.Rename(filepath.Join(oldDefaultWorkPath, file.Name()), filepath.Join(constant.DataDirPath, file.Name()))
				}
			}
		}

		err = os.RemoveAll(oldDefaultWorkPath)
		if err != nil {
			log.Errorf("Auto migration failed: %s\n", err)
		}
	}
}
