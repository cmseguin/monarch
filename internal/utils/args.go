package utils

import (
	"strconv"

	"github.com/spf13/cobra"
)

func GetStringArg(cmd *cobra.Command, cobraKey, envKey string, defaultValue string) string {
	var value string = ""

	if cmd != nil && cobraKey != "" {
		value, _ = cmd.Flags().GetString(cobraKey)
	}

	if value == "" && envKey != "" {
		value = GetEnv(envKey)
	}

	if value == "" {
		value = defaultValue
	}

	return value
}

func GetBoolArg(cmd *cobra.Command, cobraKey, envKey string, defaultValue bool) bool {
	value := GetStringArg(cmd, cobraKey, envKey, "")

	if value == "true" || value == "1" || value == "TRUE" || value == "True" {
		return true
	} else if value == "false" || value == "0" || value == "FALSE" || value == "False" {
		return false
	}

	return defaultValue
}

func GetIntArg(cmd *cobra.Command, cobraKey, envKey string, defaultValue int) int {
	var intValue int64
	var err error
	value := GetStringArg(cmd, cobraKey, envKey, "")

	if value != "" {
		// Cast the string to an int
		intValue, err = strconv.ParseInt(value, 10, 64)
	}

	if err != nil {
		return defaultValue
	}

	return int(intValue)
}
