/*
Copyright © 2025 Vicky Chhetri <vickychhetri4@gmail.com>
*/

package analyzer

import (
	"fmt"
	"path/filepath"
	"strings"
)

type FileReport struct {
	File        string
	FilePathStr string
	Folder      string
	ClassName   string
	Methods     []string
	Warnings    []SecurityWarning
}

type ModuleReport struct {
	Module string
	Files  []FileReport
}

type SecurityWarning struct {
	Level   string // HIGH, MEDIUM, LOW
	Message string
	File    string
	Line    int
	Snippet string // optional but very useful
	Rule    string // e.g. SQL_INJECTION_RAW
}

func BuildReport(module string, modulePath string) (*ModuleReport, error) {
	var files []string
	files, err := ScanPhpFiles(modulePath)

	if err != nil {
		fmt.Println("error: ", err)
		return nil, err
	}

	report := &ModuleReport{
		Module: module,
	}

	for _, file := range files {
		parsed, err := ParsePhpFiles(file)
		if err != nil || parsed == nil {
			continue
		}

		parts := strings.Split(file, string(filepath.Separator))
		var FileFolder string
		if len(parts) >= 2 {
			FileFolder = parts[len(parts)-2]
		}

		warnings := detectRawSQL(parsed.Code, file)

		report.Files = append(report.Files, FileReport{
			File:        filepath.Base(file),
			FilePathStr: file,
			Folder:      FileFolder,
			ClassName:   parsed.ClassName,
			Methods:     parsed.Methods,
			Warnings:    warnings,
		})
	}

	return report, nil
}
