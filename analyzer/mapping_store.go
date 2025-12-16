package analyzer

import "database/sql"

func CreateReport(db *sql.DB, reportType, projectPath string) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO reports (type, project_path) VALUES (?, ?)`,
		reportType, projectPath,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func SaveMapping(
	db *sql.DB,
	reportID int64,
	module string,
	controller string,
	model string,
	table string,
	controllerFile string,
	modelFile string,
) error {

	_, err := db.Exec(`
	INSERT INTO controller_model_table_map
	(report_id, module, controller, model, table_name, controller_file, model_file)
	VALUES (?, ?, ?, ?, ?, ?, ?)`,
		reportID,
		module,
		controller,
		model,
		table,
		controllerFile,
		modelFile,
	)

	return err
}
