/*
Copyright © 2025 Vicky Chhetri
*/

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/cobra"
	"github.com/vickychhetri/ci3-analyzer/analyzer"
)

var mapCmd = &cobra.Command{
	Use:   "map",
	Short: "Map Controller → Model → Tables (CI3 HMVC)",
	Long:  "Analyze CodeIgniter 3 HMVC project and store controller-model-table mapping in SQLite",
	Run: func(cmd *cobra.Command, args []string) {

		if projectPath == "" {
			fmt.Println("Project path is required")
			os.Exit(1)
		}

		fmt.Println("Mapping Project/:", projectPath)

		// --------------------------------------------------
		// Open SQLite DB
		// --------------------------------------------------
		db, err := analyzer.OpenDB()
		if err != nil {
			fmt.Println("DB error:", err)
			return
		}
		defer db.Close()

		// --------------------------------------------------
		// Create new report entry
		// --------------------------------------------------
		reportID, err := analyzer.CreateReport(db, "map", projectPath)
		if err != nil {
			fmt.Println("Failed to create report:", err)
			return
		}

		fmt.Println("Mapping Report ID:", reportID)

		// --------------------------------------------------
		// Scan HMVC modules
		// --------------------------------------------------
		modules, err := analyzer.ScanModules(projectPath)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		var wg sync.WaitGroup
		var mu sync.Mutex

		for _, module := range modules {
			wg.Add(1)

			go func(mod string) {
				defer wg.Done()

				modulePath := filepath.Join(projectPath, "application", "modules", mod)
				controllersPath := filepath.Join(modulePath, "controllers")
				modelsPath := filepath.Join(modulePath, "models")

				filepath.Walk(controllersPath, func(path string, info os.FileInfo, err error) error {
					if err != nil || info.IsDir() || filepath.Ext(path) != ".php" {
						return nil
					}

					controllerCode, err := os.ReadFile(path)
					if err != nil {
						return nil
					}

					controllerName := filepath.Base(path)
					models := analyzer.ExtractModels(string(controllerCode))

					for _, model := range models {
						modelFile := filepath.Join(modelsPath, model+".php")
						if _, err := os.Stat(modelFile); err != nil {
							continue
						}

						modelCode, err := os.ReadFile(modelFile)
						if err != nil {
							continue
						}

						tables := analyzer.ExtractTables(string(modelCode))

						for _, table := range tables {
							mu.Lock()
							_ = analyzer.SaveMapping(
								db,
								reportID,
								mod,
								controllerName,
								model,
								table,
								path,
								modelFile,
							)
							mu.Unlock()
						}
					}
					return nil
				})

			}(module)
		}

		wg.Wait()

		fmt.Println("Mapping completed successfully.")
	},
}

func init() {
	rootCmd.AddCommand(mapCmd)

	mapCmd.Flags().StringVarP(
		&projectPath,
		"project",
		"p",
		"",
		"Path to CI3 project",
	)
}
