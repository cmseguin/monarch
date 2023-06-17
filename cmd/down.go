package cmd

import (
	"errors"
	"os"

	"github.com/cmseguin/monarch/internal/types"
	"github.com/cmseguin/monarch/internal/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(downCmd)
}

var downCmd = &cobra.Command{
	Use:   "down [limitPattern]",
	Short: "Migration down",
	Run: func(cmd *cobra.Command, args []string) {
		utils.LoadEnvFile(utils.GetStringArg(cmd, "dotenvfile", "", ""))

		var limitPattern string = "*"

		if len(args) > 0 {
			limitPattern = args[0]
		}

		// TODO add support for inputting the init path
		var initPath string = "."

		migrationDir, exception := utils.GetMigrationPath(initPath)

		if exception != nil {
			exception.Terminate()
		}

		migrationObjects := []types.MigrationObject{}

		exception = utils.GetDownMigratrionObjectsFromDir(migrationDir, &migrationObjects)

		if exception != nil {
			exception.Terminate()
		}

		if len(migrationObjects) == 0 {
			utils.WrapError(errors.New("no migration files to rollback"), 0).Terminate()
		}

		// Filter out the down migrations
		db, exception := utils.InitDb(cmd)

		if exception != nil {
			exception.Explain("Error connecting to the database").SetCode(1).Terminate()
		}

		// Get the list of migrations that have already been run
		invalidMigrationKeysFromDatabase, exception := utils.GetMigrationsFromDatabase(db, false)

		if exception != nil {
			exception.Terminate()
		}

		sortedMigrations := utils.SortMigrationObjects(migrationObjects)
		sortedMigrations = utils.ReverseMigrationObjects(sortedMigrations)

		// Filter out the migrations that have already been run
		migrationObjectsToRun := utils.FilterMigrationToRun(
			limitPattern,
			sortedMigrations,
			invalidMigrationKeysFromDatabase,
		)

		if len(migrationObjectsToRun) == 0 {
			utils.WrapError(errors.New("no applied migration migrations to rollback after filtering"), 0).Terminate()
		}

		// Print the migrations that are going to be rollback
		utils.PrintStmt("The following migration will be rollback:")

		var migrationKeys []string = []string{}
		for _, migrationObject := range migrationObjectsToRun {
			migrationKeys = append(migrationKeys, migrationObject.Key)
		}

		utils.PrintOrderedList(migrationKeys)
		res := utils.AskForConfirmation("Continue?", "y")

		if !res {
			utils.PrintWarning("Aborting migration rollback")
			os.Exit(0)
		}

		// Run the migrations
		for _, migrationObject := range migrationObjectsToRun {
			fileContent, exception := utils.GetMigrationContent(migrationDir, migrationObject.File)

			if exception != nil {
				exception.Explain("Error getting migration content: " + migrationObject.File).Terminate()
			}

			// Run the migration
			exception = utils.ExecuteMigration(db, fileContent)

			if exception != nil {
				exception.Explain("Error running migration: " + migrationObject.File).Terminate()
			}

			// Update the status of the migration in the database.
			exception = utils.RollbackMigration(db, migrationObject.Key)

			if exception != nil {
				exception.Explain("Error updating the status of migration: " + migrationObject.Key).Terminate()
			}
		}

		utils.PrintSuccess("Migrations rollback successfully")
		os.Exit(0)
	},
}
