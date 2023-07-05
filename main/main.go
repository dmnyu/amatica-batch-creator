package main

import (
	amatica_batch_creator "amatica-batch-creator"
	"flag"
	"fmt"
	"os"
)

var (
	sourceDirectory  string
	targetDirectory  string
	maxSize          int64
	maxFiles         int
	sourceFilesSize  int64
	sourceFilesCount int
)

func init() {
	flag.StringVar(&sourceDirectory, "source", "", "Source directory to split into batches")
	flag.StringVar(&targetDirectory, "target", "", "Target directory to store batches")
	flag.Int64Var(&maxSize, "max-file-size", 0, "Maximum size of each batch")
	flag.IntVar(&maxFiles, "max-file-count", 0, "Maximum number of files in each batch")
}

func main() {
	flag.Parse()

	if maxSize == 0 && maxFiles == 0 {
		panic("Either max-file-size or max-file-count must be provided")
	}

	if err := amatica_batch_creator.DirectoryExists(sourceDirectory); err != nil {
		panic(err)
	}

	if err := amatica_batch_creator.DirectoryExists(targetDirectory); err != nil {
		if err := os.Mkdir(targetDirectory, 0644); err != nil {
			panic(err)
		}
	}

	//get the number of files and total size of the source directory
	var err error
	sourceFilesSize, sourceFilesCount, err = amatica_batch_creator.GetFileSizeAndCount(sourceDirectory)
	if err != nil {
		panic(err)
	}

	fmt.Println(sourceFilesSize, sourceFilesCount)
}
