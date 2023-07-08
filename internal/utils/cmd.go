package utils

import (
	"os"

	"github.com/cmseguin/khata"
	"github.com/spf13/cobra"
)

func CreateCmdHandler(handler func(cmd *cobra.Command, args []string) *khata.Khata) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		kErr := handler(cmd, args)

		if kErr != nil {
			kErr.Debug()
			os.Exit(kErr.ExitCode())
		}

		os.Exit(0)
	}
}
