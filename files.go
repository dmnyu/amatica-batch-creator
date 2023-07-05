package amatica_batch_creator

import (
	"errors"
	"os"
	"path/filepath"
)

func DirectoryExists(source string) error {
	if _, err := os.Stat(source); err == nil {
		return nil
	} else if errors.Is(err, os.ErrNotExist) {
		return os.ErrNotExist
	} else {
		return err
	}
}

func GetFileSizeAndCount(source string) (int64, int, error) {
	var filesSize int64 = 0
	var filesCount int = 0
	err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			filesSize += info.Size()
			filesCount += 1
		}
		return nil
	})

	if err != nil {
		return filesSize, filesCount, err
	}

	return filesSize, filesCount, nil
}
