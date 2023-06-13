package utils

import (
	"errors"
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
	installDir, err := FindInstallationPath()

	if err != nil {
		return "", err
	}

	migrationDir := path.Join(installDir, "migrations")

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

func FindInstallationPath() (string, error) {
	currentDir, err := os.Getwd()

	if err != nil {
		return "", err
	}

	// Check if the current directory is the installation directory
	if _, err := os.Stat(path.Join(currentDir, "migrations")); err == nil {
		return currentDir, nil
	}

	// Check each parent directory for the installation directory
	for {
		currentDir = path.Dir(currentDir)

		if currentDir == "/" {
			return "", errors.New("Could not find installation directory")
		}

		if _, err := os.Stat(path.Join(currentDir, "migrations")); err == nil {
			return currentDir, nil
		}
	}
}
