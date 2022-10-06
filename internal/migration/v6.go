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
					_ = os.Rename(filepath.Join(oldDefaultWorkPath, file.Name()), filepath.Join(constant.ConfigDirPath, file.Name()))
				}
			}
		}
		err = os.RemoveAll(oldDefaultWorkPath)
		if err != nil {
			log.Errorf("Auto migration failed: %s\n", err)
		}
	}
}
