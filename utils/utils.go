package utils

import (
	"os"
	"regexp"
)

const projectDirName = "wallet-live-coding-herman"

func GetAppRootDirectory() string {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	return string(rootPath)
}
