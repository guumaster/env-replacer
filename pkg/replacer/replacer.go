package replacer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
)

type EnvReplacer struct {
	missing       []string
	stopOnEmpty   bool
	stopOnMissing bool
	reEnv         *regexp.Regexp
}

func New(opts *Options) EnvReplacer {
	return EnvReplacer{
		missing:       []string{},
		stopOnEmpty:   opts.StopOnEmpty,
		stopOnMissing: opts.StopOnMissing,
		reEnv:         regexp.MustCompile(`\${[^}]+}`),
	}
}

func (e *EnvReplacer) Process(r io.Reader, w io.Writer) ([]string, error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		b := scanner.Bytes()
		if len(b) == 0 {
			_, _ = w.Write([]byte("\n"))
			continue
		}
		_, err := w.Write(e.reEnv.ReplaceAllFunc(b, e.replace))
		if err != nil {
			return e.missing, fmt.Errorf("can't write content: %w", err)
		}
		_, _ = w.Write([]byte("\n"))

		if err := scanner.Err(); err != nil {
			return e.missing, fmt.Errorf("error scanning: %w", err)
		}
	}
	if e.stopOnMissing && len(e.missing) > 0 {
		return e.missing, fmt.Errorf("%w: %s", MissingVariablesErr, e.missing)
	}

	return e.missing, nil
}

var MissingVariablesErr = errors.New("missing variables")

func (e EnvReplacer) GetMissingVars() []string {
	return e.missing
}

// Replace ${VAR_NAME} with its environment value
func (e *EnvReplacer) replace(b []byte) []byte {
	in := string(b)
	env := in[2 : len(in)-1]
	val, ok := os.LookupEnv(env)
	res := b
	if ok && val != "" {
		res = []byte(val)
	}
	if !ok || (e.stopOnEmpty && val == "") {
		e.missing = append(e.missing, env)
	}
	return res
}
