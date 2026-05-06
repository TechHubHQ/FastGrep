package config

import (
	"path/filepath"
	"slices"
	"strings"
)

const DefaultChunkSize = 1024 * 1024 // 1MB

var ExcludeDirs = []string{
	// VCS
	".git",

	// Dependencies
	"node_modules",
	"vendor",

	// Build output
	"dist",
	"build",
	"out",
	"bin",
	"target",
	"coverage",

	// Cache & temp
	".cache",
	"tmp",
	"temp",
	"logs",

	// IDE
	".vscode",
	".idea",

	// Infra
	".terraform",
}

func IsExcludedDir(path string) bool {
	base := strings.ToLower(filepath.Base(path))
	return slices.ContainsFunc(ExcludeDirs, func(s string) bool {
		return strings.ToLower(s) == base
	})
}
