package cmd

import (
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

		migrationDir, err := utils.GetMigrationPath(initPath)

		if err != nil {
			println("Error getting migration path")
			os.Exit(1)
		}

		migrationObjects := []types.MigrationObject{}

		err = utils.GetDownMigratrionObjectsFromDir(migrationDir, &migrationObjects)

		if err != nil {
			println("Error reading migrations directory")
			os.Exit(1)
		}

		if len(migrationObjects) == 0 {
			println("No migration files to rollback")
			os.Exit(0)
		}

		// Filter out the down migrations
		db, err := utils.InitDb(cmd)

		if err != nil {
			println(err.Error())
			os.Exit(1)
		}

		// Get the list of migrations that have already been run
		invalidMigrationKeysFromDatabase, err := utils.GetMigrationsFromDatabase(db, false)

		if err != nil {
			println("Error getting migrations")
			os.Exit(1)
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
			println("No applied migration migrations to rollback after filtering")
			os.Exit(0)
		}

		// Run the migrations
		for _, migrationObject := range migrationObjectsToRun {
			fileContent, err := utils.GetMigrationContent(migrationDir, migrationObject.File)

			if err != nil {
				println("Error getting migration content: " + migrationObject.File)
				os.Exit(1)
			}

			// Run the migration
			err = utils.ExecuteMigration(db, fileContent)

			if err != nil {
				println("Error rolling back migration: " + migrationObject.File)
				os.Exit(1)
			}

			// Update the status of the migration in the database.
			err = utils.RollbackMigration(db, migrationObject.Key)

			if err != nil {
				println("Error updating the status of migration: " + migrationObject.Key)
				os.Exit(1)
			}
		}

		println("Migrations rollback successfully")
		os.Exit(0)
	},
}
