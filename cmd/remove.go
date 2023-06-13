package cmd

import (
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/cmseguin/monarch/internal/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:   "remove [migrationName]",
	Short: "Remove a migration",
	Run: func(cmd *cobra.Command, args []string) {
		var migrationName string

		if len(args) > 0 {
			migrationName = args[0]
		}

		if migrationName == "" {
			println("migrationName is required")
			os.Exit(1)
		}

		migrationDir, err := utils.GetMigrationPath("")

		if err != nil {
			println("Error getting migration path")
			os.Exit(1)
		}

		// Get a list of files in the migrations directory
		entries, err := os.ReadDir(migrationDir)

		if err != nil {
			println("Error reading migrations directory")
			os.Exit(1)
		}

		var filesToDelete []string = []string{}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			// regex to match the migration name
			matched, _ := regexp.MatchString("[0-9]{14}-"+migrationName+".(up|down).sql", entry.Name())
			if matched {
				filesToDelete = append(filesToDelete, entry.Name())
			}
		}

		if len(filesToDelete) == 0 {
			println("No matching migrations found")
			os.Exit(1)
		}

		var response string
		fmt.Println("The following migrations will be deleted:")

		for _, file := range filesToDelete {
			fmt.Println(" - " + file)
		}

		fmt.Printf("Are you sure you want to delete %d migrations? (y/n): ", len(filesToDelete))
		fmt.Scanln(&response)

		if response != "y" {
			println("Aborting")
			os.Exit(1)
		}

		for _, file := range filesToDelete {
			err := os.Remove(path.Join(migrationDir, file))

			if err != nil {
				println("Error deleting migration")
				os.Exit(1)
			}
		}

		println("Migrations deleted successfully")
	},
}
