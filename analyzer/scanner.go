/*
Copyright © 2025 Vicky Chhetri <vickychhetri4@gmail.com>
*/

package analyzer

import (
	"io/fs"
	"os"
	"path/filepath"
)

func ScanModules(basePath string) ([]string, error) {
	modulePath := filepath.Join(basePath, "application", "modules")

	entries, err := os.ReadDir(modulePath)
	if err != nil {
		return nil, err
	}

	var modules []string
	for _, e := range entries {
		if e.IsDir() {
			modules = append(modules, e.Name())
		}
	}

	return modules, nil
}

func ScanPhpFiles(modulePath string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(modulePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".php" {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return []string{}, nil
	}
	return files, nil
}
