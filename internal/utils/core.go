package utils

import (
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
)

func ValidateMigrationName(migrationName string) bool {
	migrationName = strings.TrimSpace(migrationName)

	if migrationName == "" {
		return false
	}

	if strings.Contains(migrationName, " ") {
		return false
	}

	matched, err := regexp.MatchString("^[a-z0-9\\-]+$", migrationName)

	if !matched || err != nil {
		return false
	}

	if len(migrationName) > 255 {
		return false
	}

	return true
}

func GetMigrationPath(initPath string) (string, error) {
	// Get the current working directory
	currentDir, err := os.Getwd()

	if err != nil {
		return "", err
	}

	currentDir = path.Join(currentDir, initPath)
	migrationDir := path.Join(currentDir, "migrations")

	return migrationDir, nil
}

func GetMigrationContent(migrationDir, file string) (string, error) {
	migrationPath := path.Join(migrationDir, file)

	// Read the file
	migrationContent, err := os.ReadFile(migrationPath)

	if err != nil {
		return "", err
	}

	return string(migrationContent), nil
}

func SortUpMigrations(migrations []string) []string {
	sortedMigrations := append([]string{}, migrations...)

	sort.Strings(sortedMigrations)

	return sortedMigrations
}

func SortDownMigrations(migrations []string) []string {
	sortedMigrations := append([]string{}, migrations...)

	sort.Strings(sortedMigrations)

	// Reverse the order
	for i := len(sortedMigrations)/2 - 1; i >= 0; i-- {
		opp := len(sortedMigrations) - 1 - i
		sortedMigrations[i], sortedMigrations[opp] = sortedMigrations[opp], sortedMigrations[i]
	}

	return sortedMigrations
}
