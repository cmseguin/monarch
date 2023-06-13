package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "monarch",
	Short: "Monarch is a sql migration tool written in go",
	Long: `
                    .;;,
.,.               .,;;;;;,
;;;;;;;,,        ,;;%%%%%;;
 ';;;%%%%;;,.  ,;;%%;;%%%;;
  ';%%;;%%%;;,;;%%%%%%%;;'
     ';;%%;;%:,;%%%%%;;%%;;,
        ';;%%%,;%%%%%%%%%;;;   .  . .-. . . .-. .-. .-. . . 
           ';:%%%%%%;;%%;;;'   |\/| | | |\| |-| |(  |   |-| 
               .:::::::.       '  ' '-' ' ' ' ' ' ' '-' ' ' 
--------------- s. - s. ---------------------------------------

Monarch is a CLI library written in Go that allows you to have very simple migrations
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = "0.0.1"
	rootCmd.SetVersionTemplate("v{{.Version}}\n")
}
