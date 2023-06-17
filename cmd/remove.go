package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/cmseguin/monarch/internal/utils"
	"github.com/ryanuber/go-glob"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:   "remove [migrationName]",
	Short: "Remove a migration",
	Run: func(cmd *cobra.Command, args []string) {
		utils.LoadEnvFile(utils.GetStringArg(cmd, "dotenvfile", "", ""))

		var migrationPattern string

		if len(args) > 0 {
			migrationPattern = args[0]
		}

		if migrationPattern == "" {
			utils.WrapError(errors.New("migration name is required"), 1).Terminate()
		}

		migrationDir, exception := utils.GetMigrationPath("")

		if exception != nil {
			exception.Terminate()
		}

		// Get a list of files in the migrations directory
		entries, err := os.ReadDir(migrationDir)

		if err != nil {
			utils.WrapError(err, 1).Explain("Error reading migrations directory").Terminate()
		}

		var filesToDelete []string = []string{}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			// regex to match the migration name
			matched := glob.Glob(migrationPattern, entry.Name())
			if matched {
				filesToDelete = append(filesToDelete, entry.Name())
			}
		}

		if len(filesToDelete) == 0 {
			utils.WrapError(errors.New("no matching migrations found"), 1).Terminate()
		}

		utils.PrintStmt("The following migrations will be deleted:")

		utils.PrintUnorderedList(filesToDelete)

		res := utils.AskForConfirmation(fmt.Sprintf("Are you sure you want to delete %d migrations?", len(filesToDelete)), "n")

		if !res {
			println("Aborting removal of migrations")
			os.Exit(0)
		}

		errorFiles := []string{}
		for _, file := range filesToDelete {
			err := os.Remove(path.Join(migrationDir, file))

			if err != nil {
				errorFiles = append(errorFiles, file)
			}
		}

		if len(errorFiles) > 0 {
			utils.PrintWarning("The following migrations could not be deleted:")
			utils.PrintUnorderedList(errorFiles)
			os.Exit(1)
		}

		utils.PrintSuccess("Migrations deleted successfully")
	},
}
