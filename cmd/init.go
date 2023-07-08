package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/cmseguin/khata"
	"github.com/cmseguin/monarch/internal/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("connection", "c", "", "Database connection string")
	initCmd.Flags().StringP("driver", "d", "", "Database connection string")
	initCmd.Flags().StringP("dotenvfile", "e", "", "Env file to load")
}

var initCmd = &cobra.Command{
	Use:   "init [path] [flags]",
	Short: "Initialize monarch's migration directory & creates a table in the database to track migrations",
	Run: utils.CreateCmdHandler(func(cmd *cobra.Command, args []string) *khata.Khata {
		utils.LoadEnvFile(utils.GetStringArg(cmd, "dotenvfile", "", ""))

		// Get the current working directory
		currentDir, err := os.Getwd()

		if err != nil {
			return khata.Wrap(err).Explain("Error getting current working directory")
		}

		// Try to connect to the database
		db, kErr := utils.InitDb(cmd)

		if kErr != nil {
			return kErr.Explain("Error connecting to the database")
		}

		// Create the migrations table
		kErr = utils.CreateMigrationTable(db)

		if kErr != nil {
			return kErr.Explain("Error creating migration table")
		}

		migrationDir := path.Join(currentDir, "migrations")

		// check if the directory exists
		if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
			// create the directory
			err := os.Mkdir(migrationDir, 0755)

			if err != nil {
				return khata.Wrap(err).Explain("Error creating migrations directory")
			}
		}

		utils.PrintSuccess(fmt.Sprintf("Initialized monarch in %s", currentDir))
		return nil
	}),
}
