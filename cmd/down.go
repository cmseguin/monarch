package cmd

import (
	"os"
	"path/filepath"

	"github.com/cmseguin/monarch/internal/utils"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
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

		entries, err := os.ReadDir(migrationDir)

		if err != nil {
			println("Error reading migrations directory")
			os.Exit(1)
		}

		var filesToMigrate []string = []string{}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			// glob match the migration name
			matched, _ := filepath.Match("*-"+limitPattern+".down.sql", entry.Name())
			if matched {
				filesToMigrate = append(filesToMigrate, entry.Name())
			}
		}

		if len(filesToMigrate) == 0 {
			println("No migrations to run")
			os.Exit(0)
		}

		// Filter out the down migrations
		db, err := utils.InitDb(cmd)

		if err != nil {
			println(err.Error())
			os.Exit(1)
		}

		// Get the list of migrations that have already been run
		appliedMigrations, err := utils.GetMigrationsFromDatabase(db, true)

		if err != nil {
			println("Error getting migrations")
			os.Exit(1)
		}

		if len(appliedMigrations) == 0 {
			println("No migrations to run")
			os.Exit(0)
		}

		sortedMigrations := utils.SortDownMigrations(appliedMigrations)

		// Filter out the migrations that have not already been run
		var filteredMigrations []string = []string{}
		for _, file := range filesToMigrate {
			if slices.Contains(sortedMigrations, file) {
				filteredMigrations = append(filteredMigrations, file)
			}
		}

		if len(filteredMigrations) == 0 {
			println("No migrations to run")
			os.Exit(0)
		}

		// Run the migrations
		for _, file := range filteredMigrations {
			fileContent, err := utils.GetMigrationContent(migrationDir, file)

			if err != nil {
				println("Error getting migration content: " + file)
				os.Exit(1)
			}

			// Run the migration
			err = utils.ExecuteMigration(db, fileContent)

			if err != nil {
				println("Error running migration: " + file)
				os.Exit(1)
			}

			// Update the status of the migration in the database.
			err = utils.RollbackMigration(db, file)

			if err != nil {
				println("Error updating the status of migration: " + file)
				os.Exit(1)
			}
		}

		println("Migrations rollback successfully")
		os.Exit(0)
	},
}
