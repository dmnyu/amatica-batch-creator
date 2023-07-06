package main

import (
	amatica_batch_creator "amatica-batch-creator"
	"flag"
	"fmt"
	"path/filepath"
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
	fmt.Println("amatica-batch-creator")
	flag.Parse()

	if maxSize == 0 && maxFiles == 0 {
		panic("Either max-file-size or max-file-count must be provided")
	}

	if err := amatica_batch_creator.DirectoryExists(sourceDirectory); err != nil {
		panic(err)
	}
	amatica_batch_creator.SetSourceDirectory(sourceDirectory)

	if err := amatica_batch_creator.DirectoryExists(targetDirectory); err != nil {
		if err := amatica_batch_creator.CreateTargetDirectory(targetDirectory); err != nil {
			panic(err)
		}
	}

	//get the number of files and total size of the source directory
	fmt.Println("Counting files and calculating total size of source directory...")
	var err error
	sourceFilesSize, sourceFilesCount, err = amatica_batch_creator.GetFileSizeAndCount(sourceDirectory)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Total size of source directory: %d\n", sourceFilesSize)
	fmt.Printf("Total Count of files in source directory: %d\n", sourceFilesCount)

	//calculate the number of batch directories
	fmt.Println("Calculating number of batch directories...")
	numBatchDirs := amatica_batch_creator.GetNumBatchDirectories(sourceFilesCount, maxFiles)
	fmt.Printf("Number of batch directories: %d\n", numBatchDirs)

	//create the batch directories
	fmt.Println("Creating batch directories...")
	if err := amatica_batch_creator.CreateBatchDirectories(targetDirectory, filepath.Base(sourceDirectory), numBatchDirs); err != nil {
		panic(err)
	}

	//split the sourcefiles slice into subslices
	amatica_batch_creator.SplitSourceFileSlice()
	for i, s := range amatica_batch_creator.SourceFilesSlices {
		fmt.Println(i, len(s))
	}

	//copy the files
	amatica_batch_creator.CopyFiles()
}
