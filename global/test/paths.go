package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// GetBasePath return the project root path. The project root path can be used
// with relative paths to create absolute paths for test files.
func GetBasePath() string {
	dir, err := os.Getwd()
	if err != nil {
		panic("could not get working directory")
	}
	for _, err := os.ReadFile(filepath.Join(dir, "go.mod")); err != nil && len(dir) > 0; {
		dir = filepath.Dir(dir)
		_, err = os.ReadFile(filepath.Join(dir, "go.mod"))
	}
	if len(dir) < 2 {
		panic("No go.mod found")
	}

	return dir
}

// ReadTestFile reads a test file using the relative path of the file from project root. For example:
//
//	seqFile := ReadTestFile(t, "./service/testdata/sed.xlsx")
//
// Using relative paths from the project root helps keep the file paths more consistent throughout the test cases
// and also helps running the test cases from either root or package directories.
func ReadTestFile(t *testing.T, relPath string) []byte {
	t.Helper()

	absPath := filepath.Join(GetBasePath(), relPath)
	file, err := os.ReadFile(absPath)
	require.NoError(t, err)
	return file
}
