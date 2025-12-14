package analyzer

import (
	"regexp"
	"strings"
)

var rawSQLRegex = regexp.MustCompile(`\$this->db->query\(\s*"(.*?)"`)

func detectRawSQL(code string, filePath string) []SecurityWarning {
	var warnings []SecurityWarning

	lines := strings.Split(code, "\n")

	for i, line := range lines {
		matches := rawSQLRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			sql := matches[1]

			if strings.Contains(sql, "$") && !strings.Contains(sql, "?") {
				warnings = append(warnings, SecurityWarning{
					Level:   "HIGH",
					Message: "Possible SQL Injection: raw SQL with variable interpolation",
					File:    filePath,
					Line:    i + 1, // line numbers start from 1
					Snippet: strings.TrimSpace(line),
					Rule:    "SQL_INJECTION_RAW",
				})
			}
		}
	}

	return warnings
}
