package utils

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnvFile(filenames ...string) {
	if len(filenames) == 1 && filenames[0] == "" {
		filenames[0] = ".env"
	}

	err := godotenv.Load(filenames...)

	if err != nil {
		println(err.Error())
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}

func GetBoolEnv(key string) bool {
	val := GetEnv(key)
	if val == "" ||
		val == "0" ||
		val == "False" ||
		val == "FALSE" ||
		val == "false" {
		return false
	}

	return true
}

func GetIntEnv(key string) int {
	val := GetEnv(key)

	if val != "" {
		// Cast the string to an int
		intValue, err := strconv.ParseInt(val, 10, 64)

		if err != nil {
			return int(intValue)
		}
	}

	return 0
}
