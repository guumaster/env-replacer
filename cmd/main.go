// Package main contains code tu run env-replacer as a CLI command.
package main

import (
	"log"
	"os"
	"regexp"

	"github.com/urfave/cli/v2"

	"github.com/guumaster/env-replacer/pkg/replacer"
)

var version = "dev"

var reTrailingSlash = regexp.MustCompile(`/?$`)

func main() {
	cmd := buildCLI(&replacer.App{})
	if err := cmd.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
func buildCLI(app *replacer.App) *cli.App {

	cli.VersionFlag = &cli.BoolFlag{
		Name:  "version",
		Usage: "print the version",
	}

	return &cli.App{
		Name:        "env-replacer",
		HelpName:    "env-replacer",
		Usage:       "replace env variables on files",
		Version:     version,
		Description: "",
		HideHelp:    true,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "path", Value: "", Usage: "path to scan for templates"},
			&cli.StringFlag{Name: "files", Value: "*.tpl", Usage: "file pattern to search"},
			&cli.BoolFlag{Name: "quiet", Value: false, Usage: "hide report", Aliases: []string{"q"}},
			&cli.BoolFlag{Name: "no-recursive", Value: false, Usage: "don't scan path recursively"},
			&cli.BoolFlag{Name: "no-truncate", Value: false, Usage: "don't truncate dst files"},
			&cli.BoolFlag{Name: "no-error", Value: false, Usage: "early quit on error"},
			&cli.BoolFlag{Name: "no-missing", Value: false, Usage: "don't write files with missing variables"},
			&cli.BoolFlag{Name: "no-empty", Value: false, Usage: "fail also with variables"},
			&cli.BoolFlag{Name: "strict", Value: false, Usage: "same as --no-empty and --no-missing combined"},
		},
		Authors: []*cli.Author{
			{
				Name:  "guumaster",
				Email: "guuweb@gmail.com",
			},
		},
		Before: func(c *cli.Context) error {
			if c.Bool("version") {
				cli.ShowVersion(c)
				return cli.Exit("", 0)
			}
			return nil
		},
		Action: func(c *cli.Context) error {
			if isPiped() {
				rep := replacer.New(replacer.WithPipedOptions())
				_, err := rep.Process(os.Stdin, os.Stdout)
				if err != nil {
					return cli.Exit(err.Error(), 1)
				}

				return nil
			}

			cwd, err := os.Getwd()
			if err != nil {
				return cli.Exit("can't read current directory", 1)
			}

			if c.NArg() > 0 {
				return cli.Exit("wrong number of arguments. Put flags before path.", 1)
			}

			basePath := c.String("path")
			if basePath == "" {
				basePath = cwd
			}
			basePath = reTrailingSlash.ReplaceAllLiteralString(basePath, "")
			if _, err := os.Stat(basePath); os.IsNotExist(err) {
				return cli.Exit("path must be set to an existing directory", 1)
			}

			reporter, err := app.Run(&replacer.Options{
				BasePath:      basePath,
				FilePattern:   c.String("files"),
				ShowReport:    !c.Bool("quiet"),
				ScanRecursive: !c.Bool("no-recursive"),
				TruncateFiles: !c.Bool("no-truncate"),
				StopOnErrors:  c.Bool("no-error"),
				StopOnMissing: c.Bool("no-missing"),
				StopOnEmpty:   c.Bool("no-empty"),
				Strict:        c.Bool("strict"),
			})
			if err != nil {
				return err
			}

			reporter.Print()

			return nil
		},
	}
}

func isPiped() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	notPipe := info.Mode()&os.ModeNamedPipe == 0
	return !notPipe
}
