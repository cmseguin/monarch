package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/cmseguin/khata"
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
	Run: utils.CreateCmdHandler(func(cmd *cobra.Command, args []string) *khata.Khata {
		utils.LoadEnvFile(utils.GetStringArg(cmd, "dotenvfile", "", ""))

		var migrationPattern string

		if len(args) > 0 {
			migrationPattern = args[0]
		}

		if migrationPattern == "" {
			return khata.New("migration name is required")
		}

		migrationDir, kErr := utils.GetMigrationPath("")

		if kErr != nil {
			return kErr.Explain("Error getting migration path")
		}

		// Get a list of files in the migrations directory
		entries, err := os.ReadDir(migrationDir)

		if err != nil {
			return khata.Wrap(err).Explain("Error reading migrations directory")
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
			return khata.New("no matching migrations found")
		}

		utils.PrintStmt("The following migrations will be deleted:")
		utils.PrintUnorderedList(filesToDelete)

		res := utils.AskForConfirmation(fmt.Sprintf("Are you sure you want to delete %d migrations?", len(filesToDelete)), "n")

		if !res {
			utils.PrintWarning("Aborting removal of migrations")
			return nil
		}

		errorFiles := []string{}
		for _, file := range filesToDelete {
			err := os.Remove(path.Join(migrationDir, file))

			if err != nil {
				errorFiles = append(errorFiles, file)
			}
		}

		if len(errorFiles) > 0 {
			return khata.New("error deleting migrations").
				Explainf("Error deleting the following migrations: %v", errorFiles)
		}

		utils.PrintSuccess("Migrations deleted successfully")
		return nil
	}),
}
