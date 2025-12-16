package analyzer

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "ci3-analyzer.db")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS reports (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT,
		project_path TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS controller_model_table_map (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		report_id INTEGER,
		module TEXT,
		controller TEXT,
		model TEXT,
		table_name TEXT,
		controller_file TEXT,
		model_file TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`)

	return db, err
}
