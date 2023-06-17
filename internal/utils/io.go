package utils

import (
	"fmt"

	"github.com/ttacon/chalk"
)

func PrintStmt(stmt string) {
	fmt.Println(stmt)
}

func PrintUnorderedList(list []string) {
	for _, item := range list {
		fmt.Println(chalk.Green, " • ", chalk.Reset, item)
	}
}

func PrintOrderedList(list []string) {
	for i, item := range list {
		fmt.Println(chalk.Green, " ", i+1, ". ", chalk.Reset, item)
	}
}

func AskForConfirmation(message, defaultValue string) bool {
	var response string

	if defaultValue == "y" {
		message += " " + chalk.Bold.TextStyle("[Y/n]") + ": "
	} else if defaultValue == "n" {
		message += " " + chalk.Bold.TextStyle("[y/N]") + ": "
	} else {
		message += " " + chalk.Bold.TextStyle("[y/n]") + ": "
	}

	fmt.Println(message)
	fmt.Scanln(&response)

	if defaultValue == "y" && response == "" {
		return true
	} else if defaultValue == "n" && response == "" {
		return false
	} else {
		return response == "y"
	}
}

func PrintSuccess(message string) {
	fmt.Println(chalk.Green.Color(message))
}

func PrintInfo(message string) {
	fmt.Println(chalk.Cyan.Color(message))
}

func PrintWarning(message string) {
	fmt.Println(chalk.Yellow.Color(message))
}

func PrintErrorMessage(message string) {
	fmt.Println(chalk.Red.Color(message))
}

func PrintError(err error) {
	fmt.Println(chalk.Red.Color(err.Error()))
}

func PrintFatal(err error) {
	fmt.Println(chalk.Red.Color(err.Error()))
	panic(err)
}
