package tpl

import (
	"strings"

	"github.com/thehowl/conf"
)

// TemplateConfig are the template configuration variables that can be modified
// by using the apposite syntax anywhere in the tempalte file.
type TemplateConfig struct {
	// NoCompile is whether the template should not be considered as a
	// SimplePage on its own, and only as a file to be included
	// (or just unused.)
	NoCompile bool
	// Include is a comma-separated list of files to include relative to the
	// directory of this template itself.
	Include string
	// Template is the name of the template file, relative to the directory
	// of Hanayo. This is not meant to be modified from the TemplateConfig.
	Template string

	// Handler is the HTTP GET handler that should handle the request.
	Handler string
	// TitleBar is the title of the page in the <title> tag.
	TitleBar string
	// KyutGrill is the name of the image to use for the header.
	KyutGrill string
	// MinPrivileges are the minimum privileges an user is required to have
	// to access this page.
	MinPrivileges uint64
	// HugeHeadingRight is whether the header should be on the left or on the
	// right.
	HugeHeadingRight bool
	// AdditionalJS is a comma-separated list of JS files to include in the
	// document.
	AdditionalJS string
}

// LoadConf creates a new TemplateConfig from a string.
func LoadConf(s string) (t TemplateConfig, err error) {
	err = conf.LoadRaw(&t, []byte(s))
	return
}

// IncludedFiles returns a []string containing the files included.
// Basically, comma-splitted Include. Prefix is the optional string to prefix
// to all file names.
func (t TemplateConfig) IncludedFiles(prefix string) []string {
	x := strings.Split(t.Include, ",")
	r := make([]string, 0, len(x))
	for _, v := range x {
		if v == "" {
			continue
		}
		r = append(r, prefix+v)
	}
	return r
}
