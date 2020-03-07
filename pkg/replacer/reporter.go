package replacer

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"text/template"
)

type Reporter struct {
	writer      io.Writer
	Files       []string
	Matches     []string
	Dirs        []string
	Errors      []error
	MissesFiles []MissingVars
	vars        map[string]bool
	MissesVars  []string
}

type MissingVars struct {
	Path string
	Vars []string
}

var tpl = `
Matches: {{ len .Matches }} files.
Missing: {{ len .MissesVars }} env vars.
{{ if gt (len .MissesFiles) 0 -}}

Files:
{{- range .MissesFiles }} 
 - {{ .Path }}
   {{- range .Vars }} 
      {{ . }}
	 {{- end -}}
{{- end }}
{{- end }}
{{ if gt (len .MissesVars) 0 -}}

Vars:
{{- range .MissesVars }} 
  - {{ . }}
{{- end }}
{{ range .Errors }} 
  - {{ .Error }}
{{- end }}
{{- end }}
`

func NewReporter(w io.Writer) *Reporter {
	return &Reporter{
		writer: w,
		vars:   map[string]bool{},
	}
}

func (r Reporter) Print() {
	w := r.writer

	tpl, err := template.New("report").Parse(tpl)
	if err != nil {
		_, _ = fmt.Fprintf(w, "error parsing template: %s", err.Error())
		return
	}

	sort.Strings(r.MissesVars)
	err = tpl.Execute(w, r)
	if err != nil {
		_, _ = fmt.Fprintf(w, "error parsing template: %s", err.Error())
	}
}

func (r *Reporter) CollectError(err error, quitOnError bool) error {
	if err == nil || errors.Is(err, MissingVariablesErr) {
		return nil
	}
	r.Errors = append(r.Errors, err)
	if quitOnError {
		return err
	}
	return nil
}

func (r *Reporter) Missing(path string, vars []string) {
	r.MissesFiles = append(r.MissesFiles, MissingVars{path, vars})

	// Trick to get unique list. two lists, but one loop
	for _, v := range vars {
		if _, ok := r.vars[v]; !ok {
			r.MissesVars = append(r.MissesVars, v)
			r.vars[v] = true
		}
	}
}

func (r *Reporter) File(p string) {
	r.Files = append(r.Files, p)
}

func (r *Reporter) Dir(p string) {
	r.Dirs = append(r.Dirs, p)
}

func (r *Reporter) Match(p string) {
	r.Matches = append(r.Matches, p)
}

func (r *Reporter) Header(pattern string) {
	_, _ = fmt.Fprintf(r.writer, "Pattern: %s\n", pattern)
}
