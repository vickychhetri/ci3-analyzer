/*
Copyright © 2025 Vicky Chhetri <vickychhetri4@gmail.com>

This file contains a basic static security analyzer written in Go.
It scans PHP (CodeIgniter-style) source code to detect common security issues
such as SQL Injection, XSS, insecure file uploads, command injection,
and missing input validation.
*/

package analyzer

import (
	"regexp"
	"strings"
)

// ------------------------------------------------------------
// RAW SQL INJECTION DETECTION
// ------------------------------------------------------------

// rawSQLRegex matches CodeIgniter raw SQL queries like:
// $this->db->query("SELECT * FROM table WHERE id = $id")
var rawSQLRegex = regexp.MustCompile(`\$this->db->query\(\s*"(.*?)"`)

// detectRawSQL checks for SQL queries where variables are directly
// interpolated into the query string without parameter binding.
func detectRawSQL(code string, filePath string) []SecurityWarning {
	var warnings []SecurityWarning

	// Split the file content line-by-line
	lines := strings.Split(code, "\n")

	for i, line := range lines {
		// Try to match raw SQL query pattern
		matches := rawSQLRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			sql := matches[1]

			// If SQL contains PHP variables ($var) and no placeholders (?),
			// it is likely vulnerable to SQL Injection
			if strings.Contains(sql, "$") && !strings.Contains(sql, "?") {
				warnings = append(warnings, SecurityWarning{
					Level:   "HIGH",
					Message: "Possible SQL Injection: raw SQL with variable interpolation",
					File:    filePath,
					Line:    i + 1, // Line numbers start from 1
					Snippet: strings.TrimSpace(line),
					Rule:    "SQL_INJECTION_RAW",
				})
			}
		}
	}

	return warnings
}

// ------------------------------------------------------------
// CROSS-SITE SCRIPTING (XSS) DETECTION
// ------------------------------------------------------------

// xssRegex detects direct echoing of user input like:
// echo $_GET['name'];
var xssRegex = regexp.MustCompile(`echo\s+.*\$_(GET|POST|REQUEST)`)

// detectXSS checks if user input is echoed without escaping
// using htmlspecialchars() or htmlentities().
func detectXSS(code, filePath string) []SecurityWarning {
	var warnings []SecurityWarning
	lines := strings.Split(code, "\n")

	for i, line := range lines {
		if xssRegex.MatchString(line) &&
			!strings.Contains(line, "htmlspecialchars") &&
			!strings.Contains(line, "htmlentities") {

			warnings = append(warnings, SecurityWarning{
				Level:   "HIGH",
				Message: "Possible XSS: user input echoed without escaping",
				File:    filePath,
				Line:    i + 1,
				Snippet: strings.TrimSpace(line),
				Rule:    "XSS_UNESCAPED_OUTPUT",
			})
		}
	}
	return warnings
}

// ------------------------------------------------------------
// INSECURE FILE UPLOAD DETECTION
// ------------------------------------------------------------

// fileUploadRegex detects usage of PHP $_FILES superglobal
var fileUploadRegex = regexp.MustCompile(`\$_FILES`)

// detectFileUploadIssues checks whether file uploads are done
// without validating allowed file types or MIME types.
func detectFileUploadIssues(code, filePath string) []SecurityWarning {
	var warnings []SecurityWarning
	lines := strings.Split(code, "\n")

	for i, line := range lines {
		if fileUploadRegex.MatchString(line) &&
			!strings.Contains(code, "allowed_types") &&
			!strings.Contains(code, "mime_content_type") {

			warnings = append(warnings, SecurityWarning{
				Level:   "HIGH",
				Message: "Possible insecure file upload: missing file type validation",
				File:    filePath,
				Line:    i + 1,
				Snippet: strings.TrimSpace(line),
				Rule:    "INSECURE_FILE_UPLOAD",
			})
		}
	}
	return warnings
}

// ------------------------------------------------------------
// COMMAND INJECTION DETECTION
// ------------------------------------------------------------

// cmdInjectionRegex detects OS command execution functions
// receiving user input directly
var cmdInjectionRegex = regexp.MustCompile(`(exec|shell_exec|system|passthru)\s*\(.*\$_(GET|POST|REQUEST)`)

// detectCommandInjection flags cases where user input is passed
// directly to system-level commands.
func detectCommandInjection(code, filePath string) []SecurityWarning {
	var warnings []SecurityWarning
	lines := strings.Split(code, "\n")

	for i, line := range lines {
		if cmdInjectionRegex.MatchString(line) {
			warnings = append(warnings, SecurityWarning{
				Level:   "HIGH",
				Message: "Possible Command Injection: user input passed to OS command",
				File:    filePath,
				Line:    i + 1,
				Snippet: strings.TrimSpace(line),
				Rule:    "COMMAND_INJECTION",
			})
		}
	}
	return warnings
}

// ------------------------------------------------------------
// MISSING INPUT VALIDATION DETECTION
// ------------------------------------------------------------

// inputRegex detects direct access to PHP superglobals
var inputRegex = regexp.MustCompile(`\$_(GET|POST|REQUEST)`)

// detectMissingValidation checks if user input is used without
// any form validation, filtering, or XSS cleaning.
func detectMissingValidation(code, filePath string) []SecurityWarning {
	var warnings []SecurityWarning

	// If none of the common validation methods are found
	if !strings.Contains(code, "form_validation") &&
		!strings.Contains(code, "xss_clean") &&
		!strings.Contains(code, "filter_input") {

		lines := strings.Split(code, "\n")
		for i, line := range lines {
			if inputRegex.MatchString(line) {
				warnings = append(warnings, SecurityWarning{
					Level:   "MEDIUM",
					Message: "Missing input validation for user-supplied data",
					File:    filePath,
					Line:    i + 1,
					Snippet: strings.TrimSpace(line),
					Rule:    "MISSING_INPUT_VALIDATION",
				})
			}
		}
	}
	return warnings
}

// ------------------------------------------------------------
// MAIN SECURITY ANALYZER ENTRY POINT
// ------------------------------------------------------------

// AnalyzeSecurity runs all security checks and returns
// a consolidated list of warnings.
func AnalyzeSecurity(code, filePath string) []SecurityWarning {
	var warnings []SecurityWarning

	warnings = append(warnings, detectRawSQL(code, filePath)...)
	warnings = append(warnings, detectXSS(code, filePath)...)
	warnings = append(warnings, detectFileUploadIssues(code, filePath)...)
	warnings = append(warnings, detectCommandInjection(code, filePath)...)
	warnings = append(warnings, detectMissingValidation(code, filePath)...)

	return warnings
}
