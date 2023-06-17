package cmd

import (
	"fmt"
	"os"
	"path"

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
	Run: func(cmd *cobra.Command, args []string) {
		utils.LoadEnvFile(utils.GetStringArg(cmd, "dotenvfile", "", ""))

		// Get the current working directory
		currentDir, err := os.Getwd()

		if err != nil {
			utils.WrapError(err, 1).Explain("Error getting current working directory").Terminate()
		}

		// Try to connect to the database
		db, exception := utils.InitDb(cmd)

		if exception != nil {
			exception.Terminate()
		}

		// Create the migrations table
		exception = utils.CreateMigrationTable(db)

		if exception != nil {
			exception.Terminate()
		}

		migrationDir := path.Join(currentDir, "migrations")

		// check if the directory exists
		if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
			// create the directory
			err := os.Mkdir(migrationDir, 0755)

			if err != nil {
				utils.WrapError(err, 1).Explain("Error creating migrations directory").Terminate()
			}
		}

		utils.PrintSuccess(fmt.Sprintf("Initialized monarch in %s", currentDir))
	},
}
