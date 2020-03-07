// Package replacer contains the app with methods replace env vars.
package replacer

import (
	"errors"
	"fmt"
	"os"
	"regexp"
)

// App contains the methods to scan files and replace env vars.
type App struct{}

func (a *App) Run(opts *Options) (*Reporter, error) {
	if opts.Strict {
		opts.StopOnEmpty = true
		opts.StopOnMissing = true
	}

	pattern := fmt.Sprintf("%s/%s", opts.BasePath, opts.FilePattern)
	if opts.ScanRecursive {
		pattern = fmt.Sprintf("%s/**/%s", opts.BasePath, opts.FilePattern)
	}
	reRemoveExt := regexp.MustCompile("\\.tpl$")

	reporter := NewReporter(os.Stdout)

	if opts.ShowReport {
		reporter.Header(pattern)
	}

	walker := NewWalker(reporter, opts)

	walker.Walk(opts.BasePath, pattern, func(srcPath string) error {
		dstPath := reRemoveExt.ReplaceAllString(srcPath, "")
		dstIsNew := !fileExists(dstPath)

		// input piped to stdin
		srcFile, dstFile, err := openFiles(srcPath, dstPath, opts.TruncateFiles)
		if err != nil {
			return err
		}
		defer srcFile.Close()
		defer dstFile.Close()

		r := New(opts)
		missing, err := r.Process(srcFile, dstFile)
		if len(missing) > 0 {
			reporter.Missing(srcPath, missing)
		}
		if errors.Is(err, MissingVariablesErr) && dstIsNew {
			_ = os.Remove(dstPath)
		}
		return err
	})
	if opts.ShowReport {
		//	fmt.Println(reporter)
	}
	return reporter, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func openFiles(src, dst string, truncate bool) (*os.File, *os.File, error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return nil, nil, fmt.Errorf("can't open src file %s: %w", src, err)
	}
	srcStats, err := os.Stat(src)
	if err != nil {
		return nil, nil, fmt.Errorf("can't read from src file %s: %w", src, err)
	}

	var dstFile *os.File
	_, err = os.Stat(dst)

	exists := fileExists(dst)

	// dstFile exists
	if exists && !truncate {
		return nil, nil, fmt.Errorf("dst file is not empty %s: %w", dst, err)
	}

	dstFile, err = os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcStats.Mode())
	if err != nil {
		return nil, nil, fmt.Errorf("can't open dst file %s: %w", dst, err)
	}

	return srcFile, dstFile, nil
}
