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
			fmt.Println("Error getting current working directory")
			os.Exit(1)
		}

		// Try to connect to the database
		db, err := utils.InitDb(cmd)

		if err != nil {
			fmt.Println("Error connecting to database")
			os.Exit(1)
		}

		// Create the migrations table
		err = utils.CreateMigrationTable(db)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		migrationDir := path.Join(currentDir, "migrations")

		// check if the directory exists
		if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
			// create the directory
			err := os.Mkdir(migrationDir, 0755)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		fmt.Println("Initialized monarch in " + currentDir)
	},
}
