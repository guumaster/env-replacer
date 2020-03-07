package replacer

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar"
)

type walker struct {
	stopOnErrors  bool
	stopOnMissing bool
	reporter      *Reporter
}

func NewWalker(r *Reporter, opts *Options) *walker {
	return &walker{
		opts.StopOnErrors,
		opts.StopOnMissing,
		r,
	}
}

type processFileFunc func(path string) error

func (w *walker) Walk(basePath string, pattern string, fn processFileFunc) {

	_ = filepath.Walk(basePath, func(srcPath string, f os.FileInfo, err error) error {
		w.reporter.File(srcPath)
		if f.IsDir() {
			w.reporter.Dir(srcPath)
			return nil
		}

		match, err := doublestar.Match(pattern, srcPath)
		err = w.reporter.CollectError(err, w.stopOnErrors)
		if err != nil {
			return err
		}
		if !match {
			return nil
		}
		w.reporter.Match(srcPath)
		if errors.Is(err, os.ErrPermission) {
			err = w.reporter.CollectError(err, w.stopOnErrors)
			if err != nil {
				return err
			}
		}

		errFn := fn(srcPath)
		err = w.reporter.CollectError(errFn, w.stopOnErrors)
		if errors.Is(errFn, MissingVariablesErr) {
			return errFn
		}
		if err != nil {
			return err
		}

		return nil
	})
}
