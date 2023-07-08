package cmd

import (
	"github.com/cmseguin/khata"
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
	Run: utils.CreateCmdHandler(func(cmd *cobra.Command, args []string) *khata.Khata {
		utils.LoadEnvFile(utils.GetStringArg(cmd, "dotenvfile", "", ""))

		var limitPattern string = "*"

		if len(args) > 0 {
			limitPattern = args[0]
		}

		// TODO add support for inputting the init path
		var initPath string = "."

		migrationDir, kErr := utils.GetMigrationPath(initPath)

		if kErr != nil {
			return kErr.Explain("Error getting migration path")
		}

		migrationObjects := []types.MigrationObject{}

		kErr = utils.GetUpMigratrionObjectsFromDir(migrationDir, &migrationObjects)

		if kErr != nil {
			return kErr.Explain("Error getting migration objects")
		}

		if len(migrationObjects) == 0 {
			utils.PrintWarning("No migrations found")
			return nil
		}

		// Filter out the down migrations
		db, kErr := utils.InitDb(cmd)

		if kErr != nil {
			return kErr.Explain("Error connecting to the database")
		}

		// Get the list of migrations that have already been run
		invalidMigrationKeysFromDatabase, kErr := utils.GetMigrationsFromDatabase(db, true)

		if kErr != nil {
			return kErr.Explain("Error getting migrations from database")
		}

		sortedMigrations := utils.SortMigrationObjects(migrationObjects)

		// Filter out the migrations that have already been run
		migrationObjectsToRun := utils.FilterMigrationToRun(
			limitPattern,
			sortedMigrations,
			invalidMigrationKeysFromDatabase,
		)

		if len(migrationObjectsToRun) == 0 {
			utils.PrintWarning("no applied migration migrations to run after filtering")
			return nil
		}

		// Print the migrations that are going to be run
		utils.PrintStmt("The following migration will be run:")

		var migrationKeys []string = []string{}
		for _, migrationObject := range migrationObjectsToRun {
			migrationKeys = append(migrationKeys, migrationObject.Key)
		}

		utils.PrintOrderedList(migrationKeys)
		res := utils.AskForConfirmation("Continue?", "y")

		if !res {
			utils.PrintWarning("Aborting migration")
			return nil
		}

		// Run the migrations
		for _, migrationObject := range migrationObjectsToRun {
			fileContent, kErr := utils.GetMigrationContent(migrationDir, migrationObject.File)

			if kErr != nil {
				return kErr.Explainf("Error getting migration content: %s", migrationObject.File)
			}

			// Run the migration
			kErr = utils.ExecuteMigration(db, fileContent)

			if kErr != nil {
				return kErr.Explainf("Error running migration: %s", migrationObject.File)
			}

			// Update the status of the migration in the database.
			kErr = utils.ApplyMigration(db, migrationObject.Key)

			if kErr != nil {
				return kErr.Explainf("Error updating the status of migration: %s", migrationObject.Key)
			}
		}

		utils.PrintSuccess("Migrations run successfully")
		return nil
	}),
}
