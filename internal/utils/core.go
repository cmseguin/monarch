package utils

import (
	"os"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/cmseguin/khata"
	"github.com/cmseguin/monarch/internal/types"
	"github.com/ryanuber/go-glob"
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

func GetMigrationPath(initPath string) (string, *khata.Khata) {
	installDir, err := FindInstallationPath()

	if err != nil {
		return "", khata.Wrap(err).SetExitCode(1).Explain("Could not find installation path")
	}

	migrationDir := path.Join(installDir, "migrations")

	return migrationDir, nil
}

func GetMigrationContent(migrationDir, file string) (string, *khata.Khata) {
	migrationPath := path.Join(migrationDir, file)

	// Read the file
	migrationContent, err := os.ReadFile(migrationPath)

	if err != nil {
		return "", khata.Wrap(err).SetExitCode(1).Explain("Could not read migration file")
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

func FindInstallationPath() (string, *khata.Khata) {
	currentDir, err := os.Getwd()

	if err != nil {
		return "", khata.Wrap(err).SetExitCode(1).Explain("Could not get current directory")
	}

	// Check if the current directory is the installation directory
	if _, err := os.Stat(path.Join(currentDir, "migrations")); err == nil {
		return currentDir, nil
	}

	// Check each parent directory for the installation directory
	for {
		currentDir = path.Dir(currentDir)

		if currentDir == "/" {
			return "", khata.New("Could not find installation directory").SetExitCode(1)
		}

		if _, err := os.Stat(path.Join(currentDir, "migrations")); err == nil {
			return currentDir, nil
		}
	}
}

func GetDownMigratrionObjectsFromDir(
	dirname string,
	migrationObjects *[]types.MigrationObject,
) *khata.Khata {
	entries, err := os.ReadDir(dirname)

	if err != nil {
		return khata.Wrap(err).SetExitCode(1).Explain("Could not read directory")
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if entry.Name()[len(entry.Name())-9:] == ".down.sql" {
			*migrationObjects = append(*migrationObjects, types.MigrationObject{
				Key:  entry.Name()[0 : len(entry.Name())-9],
				File: entry.Name(),
			})
		}
	}

	return nil
}

func GetUpMigratrionObjectsFromDir(
	dirname string,
	migrationObjects *[]types.MigrationObject,
) *khata.Khata {
	entries, err := os.ReadDir(dirname)

	if err != nil {
		return khata.Wrap(err).SetExitCode(1).Explain("Could not read directory")
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if entry.Name()[len(entry.Name())-7:] == ".up.sql" {
			*migrationObjects = append(*migrationObjects, types.MigrationObject{
				Key:  entry.Name()[0 : len(entry.Name())-7],
				File: entry.Name(),
			})
		}
	}

	return nil
}

func SortMigrationObjects(migrationObjects []types.MigrationObject) []types.MigrationObject {
	sortedMigrationObjects := append([]types.MigrationObject{}, migrationObjects...)

	sort.Slice(sortedMigrationObjects, func(i, j int) bool {
		return sortedMigrationObjects[i].Key < sortedMigrationObjects[j].Key
	})

	return sortedMigrationObjects
}

func ReverseMigrationObjects(migrationObjects []types.MigrationObject) []types.MigrationObject {
	reversedMigrationObjects := append([]types.MigrationObject{}, migrationObjects...)

	for i := len(reversedMigrationObjects)/2 - 1; i >= 0; i-- {
		opp := len(reversedMigrationObjects) - 1 - i
		reversedMigrationObjects[i], reversedMigrationObjects[opp] = reversedMigrationObjects[opp], reversedMigrationObjects[i]
	}

	return reversedMigrationObjects
}

func FilterMigrationToRun(
	limitPattern string,
	migrationObjects []types.MigrationObject,
	invalidMigrationKeysFromDatabase []string,
) []types.MigrationObject {
	migrationObjectsToRun := []types.MigrationObject{}

	for _, migrationObject := range migrationObjects {
		foundIndex := FindIndexInString(invalidMigrationKeysFromDatabase, func(value string, index int) bool {
			return value == migrationObject.Key
		})

		if foundIndex == -1 {
			migrationObjectsToRun = append(migrationObjectsToRun, migrationObject)
		}

		// If limit pattern is met, stop migrating
		if limitPattern != "" && !glob.Glob(limitPattern, migrationObject.Key) {
			break
		}
	}

	return migrationObjectsToRun
}
