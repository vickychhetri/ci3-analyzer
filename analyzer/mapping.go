/*
Copyright © 2025 Vicky Chhetri

CI3 Controller → Model → Table mapping analyzer
Supports Query Builder + Raw SQL
*/

package analyzer

import (
	"regexp"
)

// ------------------------------------------------------------
// REGEX DEFINITIONS (CI3 REAL-WORLD SUPPORT)
// ------------------------------------------------------------

// $this->load->model('User_model');
var loadModelRegex = regexp.MustCompile(
	`\$this->load->model\(\s*['"]([^'"]+)['"]`,
)

// Query Builder table usage
// from(), get(), get_where(), insert(), update(), delete(), join()
var qbTableRegex = regexp.MustCompile(
	`->(from|get|get_where|insert|update|delete|join)\(\s*['"]([^'"]+)['"]`,
)

// Raw SQL queries
// $this->db->query("SELECT ...");
var rawSQLRegexQuery = regexp.MustCompile(
	`\$this->db->query\(\s*["']([\s\S]*?)["']\s*\)`,
)

// SQL table extraction
// FROM table | JOIN table | INSERT INTO table | UPDATE table
var sqlTableRegex = regexp.MustCompile(
	`(?i)\b(from|join|into|update)\s+([a-zA-Z0-9_]+)`,
)

// ------------------------------------------------------------
// DATA STRUCTURE (optional future use)
// ------------------------------------------------------------

type Mapping struct {
	Controller string
	Model      string
	Table      string
	File       string
	Line       int
}

// ------------------------------------------------------------
// CONTROLLER → MODELS
// ------------------------------------------------------------

// ExtractModels finds all models loaded in a controller
func ExtractModels(code string) []string {
	var models []string

	matches := loadModelRegex.FindAllStringSubmatch(code, -1)
	for _, m := range matches {
		models = append(models, m[1])
	}

	return unique(models)
}

// ------------------------------------------------------------
// MODEL → TABLES (MAIN LOGIC)
// ------------------------------------------------------------

// ExtractTables finds all tables used in a model file
func ExtractTables(code string) []string {
	var tables []string

	// --------------------------------------------
	// 1. Query Builder (CI3 style)
	// --------------------------------------------
	qbMatches := qbTableRegex.FindAllStringSubmatch(code, -1)
	for _, m := range qbMatches {
		// m[2] = table name
		tables = append(tables, m[2])
	}

	// --------------------------------------------
	// 2. Raw SQL ($this->db->query)
	// --------------------------------------------
	rawMatches := rawSQLRegexQuery.FindAllStringSubmatch(code, -1)
	for _, m := range rawMatches {
		sql := m[1]

		sqlMatches := sqlTableRegex.FindAllStringSubmatch(sql, -1)
		for _, sm := range sqlMatches {
			// sm[2] = table name
			tables = append(tables, sm[2])
		}
	}

	sqlMatchesFrom := sqlTableRegex.FindAllStringSubmatch(code, -1)
	for _, sm := range sqlMatchesFrom {
		// sm[2] = table name
		tables = append(tables, sm[2])
	}
	return unique(tables)
}

// ------------------------------------------------------------
// UTILITY
// ------------------------------------------------------------

func unique(items []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, v := range items {
		if v == "" {
			continue
		}
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}
