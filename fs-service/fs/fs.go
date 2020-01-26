package fs

import (
	"os"
	"path/filepath"
)

var root string = "storage"

type FileSystem interface {
	CreateFolder(path string, name string)
	DeleteFolder(dir string)
	CreateFile()
	DeleteFile()
}

func CreateFolder(userID string) error {
	dir := filepath.Join(root, userID)
	if checkDirIfNotExists(root) {
		os.Mkdir(root, os.ModePerm)
	}
	if checkDirIfNotExists(dir) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteFolder(dir string) {
	if !checkDirIfNotExists(dir) {
		err := os.Remove(dir)
		if err != nil {
			panic(err)
		}
	}
}

func CreateFile() {

}

func DeleteFile() {
}

func checkDirIfNotExists(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return true
	}
	return false
}
