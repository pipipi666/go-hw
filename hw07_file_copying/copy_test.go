package main

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func deleteTestFile(fileName string) {
	os.Remove(fileName)
}

func TestCopy(t *testing.T) {
	testFile, err := os.CreateTemp("", "out.txt")
	if err != nil {
		fmt.Println("Create tpm file error")
	}

	testFileName := testFile.Name()

	defer deleteTestFile(testFileName)

	t.Run("Simple test", func(t *testing.T) {
		err := Copy("testdata/input.txt", testFileName, 0, 100)

		require.NoError(t, err)
	})

	t.Run("EOF test", func(t *testing.T) {
		err := Copy("testdata/input.txt", testFileName, 0, 8000)

		require.NoError(t, err)
	})

	t.Run("Offset exceeds file size test", func(t *testing.T) {
		err := Copy("testdata/input.txt", testFileName, 8000, 100)

		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual err - %v", err)
	})

	t.Run("Unsupported file test", func(t *testing.T) {
		err := Copy("\000", testFileName, 0, 100)

		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual err - %v", err)
	})
}
