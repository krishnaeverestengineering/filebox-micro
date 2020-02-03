package fs

import (
	"fmt"
	"io/ioutil"
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

func CreateFolder(path string, data UserFile) error {
	if !exists(root) {
		os.Mkdir(root, os.ModePerm)
	}
	if !exists(filepath.Join(root, data.UserID)) {
		os.Mkdir(filepath.Join(root, data.UserID), os.ModePerm)
	}
	if data.Type == Folder {
		createFolder(data.UserID, path, data.FileName)
	} else if data.Type == TextFile {
		createFile(data.UserID, path, data.FileName, data.Extension)
	}
	return nil
}

func createFolder(userID string, path string, name string) error {
	dir := filepath.Join(root, userID, path, name)
	if !exists(dir) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Folder already exists")
}

func DeleteFolder(dir string, name string, uid string) error {
	fullPath := filepath.Join(root, uid, dir)
	fmt.Println(exists(fullPath))
	if exists(fullPath) {
		err := os.RemoveAll(fullPath)
		return err
	}
	return fmt.Errorf("File or Folder not found")
}

func ReadFile(name string, userID string) (interface{}, error) {
	path := filepath.Join(root, userID, name)
	if !exists(path) {
		return nil, fmt.Errorf("File not found")
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	text := string(data)
	return text, nil
}

func EditTextFile(name string, content string, userID string) error {
	fullPath := filepath.Join(root, userID, name)
	if !exists(fullPath) {
		return fmt.Errorf("File not found")
	}
	return ioutil.WriteFile(fullPath, []byte(content), os.ModePerm)
}

func createFile(userID string, path string, name string, ext string) {
	dir := filepath.Join(root, userID, path, name+ext)
	file, err := os.Create(dir)
	fmt.Println(file, err)
}

func DeleteFile() {
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
