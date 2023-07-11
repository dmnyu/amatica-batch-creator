package amatica_batch_creator

import "testing"

var source = "test/simple-data"

func TestDirectoriesExists(t *testing.T) {
	t.Run("Test that directory exists", func(t *testing.T) {
		if err := DirectoryExists(source); err != nil {
			t.Errorf("DirectoryExists() error = %v", err)
		}
	})
}

func TestGetFileSizeAndCount(t *testing.T) {
	t.Run("Test that file size and count is correct", func(t *testing.T) {
		sourceFilesSize, sourceFilesCount, err := GetFileSizeAndCount(source)
		if err != nil {
			t.Errorf("GetFileSizeAndCount() error = %v", err)
		}
		t.Log(sourceFilesSize, sourceFilesCount)
		if sourceFilesSize != 7 {
			t.Errorf("GetFileSizeAndCount() sourceFilesSize = %v, want %v", sourceFilesSize, 7)
		}

		if sourceFilesCount != 2 {
			t.Errorf("GetFileSizeAndCount() sourceFilesCount = %v, want %v", sourceFilesCount, 2)
		}

	})
}
