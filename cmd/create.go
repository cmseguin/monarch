package cmd

import (
	"os"
	"path"
	"time"

	"github.com/cmseguin/khata"
	"github.com/cmseguin/monarch/internal/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create [migrationName]",
	Short: "Create a migration",
	Run: utils.CreateCmdHandler(func(cmd *cobra.Command, args []string) *khata.Khata {
		utils.LoadEnvFile(utils.GetStringArg(cmd, "dotenvfile", "", ""))

		var migrationName string

		if len(args) > 0 {
			migrationName = args[0]
		}

		if migrationName == "" {
			return khata.New("migration name is required")
		}

		// Get a timestamp for the migration
		datestamp := time.Now().Format("20060102150405")

		// validate the migration name
		if !utils.ValidateMigrationName(migrationName) {
			return khata.New("invalid migration name")
		}

		// Create the migration file
		migrationUpFile := datestamp + "-" + migrationName + ".up.sql"
		migrationDownFile := datestamp + "-" + migrationName + ".down.sql"

		migrationPath, kErr := utils.GetMigrationPath("")

		if kErr != nil {
			return kErr.Explain("Error getting migration path")
		}

		migrationUpPath := path.Join(migrationPath, migrationUpFile)
		migrationDownPath := path.Join(migrationPath, migrationDownFile)

		_, err := os.Stat(migrationUpPath)

		if err == nil {
			utils.PrintWarning("Migration up file already exists")
		} else if os.IsNotExist(err) {
			_, err := os.Create(migrationUpPath)
			if err != nil {
				return khata.New("error creating migration up file")
			}
		}

		_, err = os.Stat(migrationDownPath)

		if err == nil {
			utils.PrintWarning("Migration down file already exists")
		} else if os.IsNotExist(err) {
			_, err := os.Create(migrationDownPath)

			if err != nil {
				return khata.New("error creating migration down file")
			}
		}

		utils.PrintSuccess("Migration files created successfully")
		return nil
	}),
}
