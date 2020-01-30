package fs

import (
	"fmt"
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

func CreateFolder(userID string, path string, name string) error {
	dir := filepath.Join(root, userID, path, name)
	if !exists(root) {
		os.Mkdir(root, os.ModePerm)
	}
	if !exists(filepath.Join(root, userID)) {
		os.Mkdir(filepath.Join(root, userID), os.ModePerm)
	}
	if !exists(dir) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Folder already exists")
}

func DeleteFolder(dir string) {
	if exists(dir) {
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

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
