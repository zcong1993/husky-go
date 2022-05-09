package main

import (
	"embed"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/urfave/cli/v2"
)

var (
	version = "master"
	commit  = ""
	date    = ""
	builtBy = ""
)

const (
	defaultDir = ".husky"
)

//go:embed husky.sh
var fs embed.FS

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println(buildVersion(version, commit, date, builtBy)) //nolint
	}

	app := &cli.App{
		Name:    "husky-go",
		Usage:   "husky git hooks manager in go",
		Version: version,
		Commands: []*cli.Command{
			{
				Name:      "install",
				Usage:     "setup git hooks",
				ArgsUsage: "[dir] (default: .husky)",
				Action:    install,
			},
			{
				Name:   "uninstall",
				Usage:  "unset git hooks",
				Action: uninstall,
			},
			{
				Name:      "add",
				Usage:     "add cmd to hooks file",
				ArgsUsage: "[file] [cmd]",
				Action: func(context *cli.Context) error {
					if context.Args().Len() != 2 {
						return errors.New("invalid args, see --help")
					}

					file, cmd := context.Args().Get(0), context.Args().Get(1)
					return add(file, cmd)
				},
			},
			{
				Name:      "set",
				Usage:     "create(replace) hooks file",
				ArgsUsage: "[file] [cmd]",
				Action: func(context *cli.Context) error {
					if context.Args().Len() != 2 {
						return errors.New("invalid args, see --help")
					}

					file, cmd := context.Args().Get(0), context.Args().Get(1)
					return set(file, cmd)
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func install(ctx *cli.Context) error {
	if os.Getenv("HUSKY") == "0" {
		fmt.Println("HUSKY env variable is set to 0, skipping install")
		return nil
	}

	dir := ctx.Args().First()
	if dir == "" {
		dir = defaultDir
	}

	// 1. check git rev-parse
	if err := exec.Command("git", "rev-parse").Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() != 0 {
				return nil
			}
		}
		return err
	}

	// Custom dir help
	url := "https://typicode.github.io/husky/#/?id=custom-directory"

	// 2. Ensure that we're not trying to install outside of cwd
	cwd := mustCwd()
	if !strings.HasPrefix(path.Join(cwd, dir), cwd) {
		return errors.New(fmt.Sprintf(".. not allowed (see %s)", url))
	}

	// 3. Ensure that cwd is git top level
	if !exists(".git") {
		return errors.New(fmt.Sprintf(".git can't be found (see %s)", url))
	}

	// Create .husky/_
	err := os.MkdirAll(path.Join(dir, "_"), 0o777)
	if err != nil {
		return fmt.Errorf("mkdir error %v", err)
	}

	// Create .husky/_/.gitignore
	err = os.WriteFile(path.Join(dir, "_/.gitignore"), []byte("*"), 0o600)
	if err != nil {
		return fmt.Errorf("create _/.gitignore error %v", err)
	}

	// Copy husky.sh to .husky/_/husky.sh
	f, err := fs.ReadFile("husky.sh")
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(dir, "_/husky.sh"), f, 0o600)
	if err != nil {
		return fmt.Errorf("create _/.gitignore error %v", err)
	}

	// Configure repo
	err = exec.Command("git", "config", "core.hooksPath", dir).Run()
	if err != nil {
		return fmt.Errorf("git hooks failed to install %v", err)
	}

	fmt.Println("Git hooks installed")

	return nil
}

func uninstall(ctx *cli.Context) error {
	return exec.Command("git", "config", "--unset", "core.hooksPath").Run()
}

func set(file, cmd string) error {
	d := path.Dir(file)
	if !exists(d) {
		return fmt.Errorf("can't create hook, %s directory doesn't exist (try running husky install)", d)
	}

	err := os.WriteFile(file, []byte(fmt.Sprintf(`#!/usr/bin/env sh
. "$(dirname -- "$0")/_/husky.sh"

%s
`, cmd)), 0o0755)

	if err != nil {
		return err
	}

	fmt.Printf("created %s\n", file)

	return nil
}

func add(file, cmd string) error {
	if !exists(file) {
		return set(file, cmd)
	}

	// append file
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0o0755)

	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(cmd + "\n"); err != nil {
		return err
	}

	fmt.Printf("updated %s\n", file)

	return nil
}

func mustCwd() string {
	cwd, _ := os.Getwd()
	return cwd
}

func exists(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
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
