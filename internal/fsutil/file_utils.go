package fsutil

import (
	"fastgrep/internal/config"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

/*
GetGlobs - gets the matching filenames for a given incomplete
filename or * for full directory search
*/
func GetGlobs(args []string) ([]string, error) {
	files := make([]string, 0)

	for _, arg := range args {
		fileMatches, err := filepath.Glob(arg)
		if err != nil {
			return nil, err
		}

		if len(fileMatches) > 0 {
			files = append(files, fileMatches...)
		} else {
			files = append(files, arg)
		}
	}
	return files, nil
}

func FilterFiles(files []string, recursive bool) (int, []string) {
	var numValidFiles = 0
	var validFiles = make([]string, 0)

	for _, file := range files {
		fileInfo, err := os.Stat(file)
		if err != nil {
			continue
		}

		if fileInfo.IsDir() {
			if config.IsExcludedDir(file) {
				continue
			}
			if recursive {
				err := filepath.WalkDir(file, func(path string, d fs.DirEntry, err error) error {
					if err != nil {
						return err
					}

					if d.IsDir() {
						if config.IsExcludedDir(path) {
							return filepath.SkipDir
						}
						return nil
					}

					// Skip binary executables even in subdirectories
					if filepath.Ext(path) == ".exe" {
						return nil
					}
					validFiles = append(validFiles, path)
					numValidFiles++
					return nil
				})
				if err != nil {
					fmt.Printf("Error walking directory %s: %v\n", file, err)
				}
			} else {
				fmt.Printf("%s is a Directory... skipping (use -r to search recursively)\n", file)
			}
			//Continue to the next file in the list, don't add the directory itself
			continue
		}

		// Handle individual files passed as arguments
		if filepath.Ext(file) == ".exe" {
			fmt.Printf("%s is an executable Binary file... skipping\n", file)
			continue
		}

		numValidFiles++
		validFiles = append(validFiles, file)
	}

	return numValidFiles, validFiles
}

func GetFileSize(file string) (int, error) {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return 0, fmt.Errorf("error calculating file size %s", err.Error())
	}

	return int(fileInfo.Size()), nil
}
