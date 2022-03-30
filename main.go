package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	version = "master"
	commit  = ""
	date    = ""
	builtBy = ""
)

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println(buildVersion(version, commit, date, builtBy)) //nolint
	}

	app := &cli.App{
		Name:  "boom",
		Usage: "make an explosive entrance",
		Action: func(c *cli.Context) error {
			fmt.Println("boom! I say!") //nolint

			return nil
		},
		Version: version,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func buildVersion(version, commit, date, builtBy string) string {
	result := version

	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}

	if date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, date)
	}

	if builtBy != "" {
		result = fmt.Sprintf("%s\nbuilt by: %s", result, builtBy)
	}

	return result
}
