package cmd

import (
	"os"
	"path"
	"time"

	"github.com/cmseguin/monarch/internal/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create [migrationName]",
	Short: "Create a migration",
	Run: func(cmd *cobra.Command, args []string) {
		utils.LoadEnvFile(utils.GetStringArg(cmd, "dotenvfile", "", ""))

		var migrationName string

		if len(args) > 0 {
			migrationName = args[0]
		}

		if migrationName == "" {
			println("migrationName is required")
			os.Exit(1)
		}

		// Get a timestamp for the migration
		datestamp := time.Now().Format("20060102150405")

		// validate the migration name
		if !utils.ValidateMigrationName(migrationName) {
			println("Invalid migration name")
			os.Exit(1)
		}

		// Create the migration file
		migrationUpFile := datestamp + "-" + migrationName + ".up.sql"
		migrationDownFile := datestamp + "-" + migrationName + ".down.sql"

		migrationPath, err := utils.GetMigrationPath("")

		if err != nil {
			println("Error getting migration path")
			os.Exit(1)
		}

		migrationUpPath := path.Join(migrationPath, migrationUpFile)
		migrationDownPath := path.Join(migrationPath, migrationDownFile)

		_, err = os.Stat(migrationUpPath)

		if err == nil {
			println("Migration up file already exists")
		} else if os.IsNotExist(err) {
			_, err := os.Create(migrationUpPath)

			if err != nil {
				println("Error creating migration up file")
				os.Exit(1)
			}
		}

		_, err = os.Stat(migrationDownPath)

		if err == nil {
			println("Migration down file already exists")
		} else if os.IsNotExist(err) {
			_, err := os.Create(migrationDownPath)

			if err != nil {
				println("Error creating migration down file")
				os.Exit(1)
			}
		}

		println("Migration files created successfully")
	},
}
