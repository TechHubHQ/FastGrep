package main

import (
	"fastgrep/internal/orchestrator"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func createLargeTestFile(t testing.TB, sizeMB int) string {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "large_test_file.txt")
	f, err := os.Create(filePath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	line := "This is a test line that will be repeated many times to reach the desired file size. SearchPattern is here.\n"
	iterations := (sizeMB * 1024 * 1024) / len(line)
	for range iterations {
		f.WriteString(line)
	}
	return filePath
}

func BenchmarkEndToEnd(b *testing.B) {
	// Create a 50MB file
	filePath := createLargeTestFile(b, 50)
	searchString := "SearchPattern"

	// Redirect stdout to avoid cluttering benchmark results
	oldStdout := os.Stdout
	devNull, _ := os.Open(os.DevNull)
	os.Stdout = devNull
	defer func() {
		os.Stdout = oldStdout
		devNull.Close()
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := orchestrator.Execute(searchString, []string{filePath}, false, false)
		if err != nil {
			b.Fatalf("Execute failed: %v", err)
		}
	}
}

func BenchmarkEndToEndMultipleFiles(b *testing.B) {
	// Create 5 files of 10MB each
	var filePaths []string
	tmpDir := b.TempDir()
	for i := 0; i < 5; i++ {
		filePath := filepath.Join(tmpDir, fmt.Sprintf("file_%d.txt", i))
		f, _ := os.Create(filePath)
		line := "This is a test line for multi-file search. SearchPattern is here.\n"
		iterations := (10 * 1024 * 1024) / len(line)
		for j := 0; j < iterations; j++ {
			f.WriteString(line)
		}
		f.Close()
		filePaths = append(filePaths, filePath)
	}
	searchString := "SearchPattern"

	oldStdout := os.Stdout
	devNull, _ := os.Open(os.DevNull)
	os.Stdout = devNull
	defer func() {
		os.Stdout = oldStdout
		devNull.Close()
	}()

	for b.Loop() {
		err := orchestrator.Execute(searchString, filePaths, false, false)
		if err != nil {
			b.Fatalf("Execute failed: %v", err)
		}
	}
}
