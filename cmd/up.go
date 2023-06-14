package cmd

import (
	"errors"
	"os"

	"github.com/cmseguin/monarch/internal/types"
	"github.com/cmseguin/monarch/internal/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(upCmd)
}

var upCmd = &cobra.Command{
	Use:   "up [limitPattern]",
	Short: "Migration up",
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

		exception = utils.GetUpMigratrionObjectsFromDir(migrationDir, &migrationObjects)

		if exception != nil {
			exception.Terminate()
		}

		if len(migrationObjects) == 0 {
			utils.WrapError(errors.New("no migration files to run"), 0).Terminate()
		}

		// Filter out the down migrations
		db, exception := utils.InitDb(cmd)

		if exception != nil {
			exception.Explain("Error connecting to the database").Terminate()
		}

		// Get the list of migrations that have already been run
		invalidMigrationKeysFromDatabase, exception := utils.GetMigrationsFromDatabase(db, true)

		if exception != nil {
			exception.Terminate()
		}

		sortedMigrations := utils.SortMigrationObjects(migrationObjects)

		// Filter out the migrations that have already been run
		migrationObjectsToRun := utils.FilterMigrationToRun(
			limitPattern,
			sortedMigrations,
			invalidMigrationKeysFromDatabase,
		)

		if len(migrationObjectsToRun) == 0 {
			utils.WrapError(errors.New("no applied migration migrations to run after filtering"), 0).Terminate()
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
			exception = utils.ApplyMigration(db, migrationObject.Key)

			if exception != nil {
				exception.Explain("Error updating the status of migration: " + migrationObject.Key).Terminate()
			}
		}

		println("Migrations run successfully")
		os.Exit(0)
	},
}
