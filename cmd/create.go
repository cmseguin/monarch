package cmd

import (
	"errors"
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
			utils.WrapError(errors.New("migration name is required"), 1).Terminate()
		}

		// Get a timestamp for the migration
		datestamp := time.Now().Format("20060102150405")

		// validate the migration name
		if !utils.ValidateMigrationName(migrationName) {
			utils.PrintWarning("Invalid migration name")
			os.Exit(1)
		}

		// Create the migration file
		migrationUpFile := datestamp + "-" + migrationName + ".up.sql"
		migrationDownFile := datestamp + "-" + migrationName + ".down.sql"

		migrationPath, exception := utils.GetMigrationPath("")

		if exception != nil {
			exception.Terminate()
		}

		migrationUpPath := path.Join(migrationPath, migrationUpFile)
		migrationDownPath := path.Join(migrationPath, migrationDownFile)

		_, err := os.Stat(migrationUpPath)

		if err == nil {
			utils.PrintWarning("Migration up file already exists")
		} else if os.IsNotExist(err) {
			_, err := os.Create(migrationUpPath)
			if err != nil {
				utils.WrapError(err, 1).Explain("Error creating migration up file").Terminate()
			}
		}

		_, err = os.Stat(migrationDownPath)

		if err == nil {
			utils.PrintWarning("Migration down file already exists")
		} else if os.IsNotExist(err) {
			_, err := os.Create(migrationDownPath)

			if err != nil {
				utils.WrapError(err, 1).Explain("Error creating migration down file").Terminate()
			}
		}

		utils.PrintSuccess("Migration files created successfully")
	},
}
