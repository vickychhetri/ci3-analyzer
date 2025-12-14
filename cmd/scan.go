/*
Copyright © 2025 Vicky Chhetri <vickychhetri4@gmail.com>
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

var projectPath string
var outputHTML bool

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan CI3 HMVC project",
	Long:  "Scan CodeIgniter 3 HMVC project and extract modules, classes, and methods",
	Run: func(cmd *cobra.Command, args []string) {
		if projectPath == "" {
			fmt.Println("Project path is required")
			os.Exit(1)
		}

		fmt.Println("Scanning Project/: ", projectPath)

		modules, err := analyzer.ScanModules(projectPath)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		fmt.Println("Modules Found:")

		chReport := make(chan *analyzer.ModuleReport)
		var m sync.Mutex
		var w sync.WaitGroup
		var reports []analyzer.ModuleReport

		for _, m := range modules {
			w.Add(1)
			go func(module string) {
				defer w.Done()
				fmt.Println(" -", module)
				modulePath := filepath.Join(projectPath, "application", "modules", module)

				report, err := analyzer.BuildReport(module, modulePath)
				if err != nil {
					fmt.Println("error : ", err)
					return
				}
				chReport <- report

			}(m)

		}

		go func() {
			w.Wait()
			close(chReport)
		}()

		for rep := range chReport {
			m.Lock()
			reports = append(reports, *rep)
			m.Unlock()

			fmt.Println("Module: ", m)
			for _, f := range rep.Files {
				fmt.Printf(" - %s, (%d methods)\n", f.ClassName, len(f.Methods))
			}
		}

		if outputHTML {
			err := analyzer.GenerateHTMLReport("ci3-reports.html", reports)
			if err != nil {
				fmt.Println("HTML Generation Failed: ", err)
				return
			}

			fmt.Println("HTML Report Generated:  ci3-reprot.html")
		}

	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVarP(
		&projectPath,
		"path",
		"p",
		"",
		"Path to CI3 project",
	)

	scanCmd.Flags().BoolVar(&outputHTML, "html", false, "Generate HTML report")

}
