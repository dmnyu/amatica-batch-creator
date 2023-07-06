package amatica_batch_creator

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const permsMode = 0777

var TargetBatchDirs = []string{}
var SourceDirectory string
var Sourcefiles = []string{}
var SourceFilesSlices [][]string

func SetSourceDirectory(source string) {
	SourceDirectory = source
}

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
			Sourcefiles = append(Sourcefiles, info.Name())
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

func CreateTargetDirectory(targetDirectory string) error {
	if err := os.Mkdir(targetDirectory, permsMode); err != nil {
		return err
	}
	return nil
}
func GetNumBatchDirectories(sourceFilesCount int, maxFiles int) int {
	numBatchDirs := sourceFilesCount / maxFiles
	if maxFiles*numBatchDirs < sourceFilesCount {
		numBatchDirs += 1
	}
	return numBatchDirs
}

func CreateBatchDirectories(targetDirectory string, baseName string, numBatchDirs int) error {
	for i := 0; i < numBatchDirs; i++ {
		batchDir := filepath.Join(targetDirectory, fmt.Sprintf("%s-batch-%d", baseName, i+1))
		if err := os.Mkdir(batchDir, permsMode); err != nil {
			return err
		}
		TargetBatchDirs = append(TargetBatchDirs, batchDir)
	}
	return nil
}

func SplitSourceFileSlice() {
	SourceFilesSlices = [][]string{}

	chunkSize := (len(Sourcefiles) + len(TargetBatchDirs) - 1) / len(TargetBatchDirs)

	for i := 0; i < len(Sourcefiles); i += chunkSize {
		end := i + chunkSize

		if end > len(Sourcefiles) {
			end = len(Sourcefiles)
		}

		SourceFilesSlices = append(SourceFilesSlices, Sourcefiles[i:end])
	}
}

func CopyFiles() {
	for i, sourceFileSlice := range SourceFilesSlices {
		fmt.Println("copying files in batch", i+1)
		for _, sourceFile := range sourceFileSlice {
			sourcePath := filepath.Join(SourceDirectory, sourceFile)
			targetPath := filepath.Join(TargetBatchDirs[i], sourceFile)
			fmt.Println(sourcePath + " -> " + targetPath)
			copyFile(sourcePath, targetPath)
		}
	}
}

func copyFile(src string, dst string) error {
	sfi, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return err
		}
	}

	if err := copyFileContents(src, dst); err != nil {
		fmt.Println("error copying file", err.Error())
	}
	return nil
}

func copyFileContents(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	err = out.Sync()
	return err
}
